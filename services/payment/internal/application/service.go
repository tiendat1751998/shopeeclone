package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopee-clone/shopee/services/payment/internal/config"
	"github.com/shopee-clone/shopee/services/payment/internal/domain"
	"github.com/shopee-clone/shopee/services/payment/internal/infrastructure/kafka"
	"github.com/shopee-clone/shopee/services/payment/internal/infrastructure/mysql"
	redisinfra "github.com/shopee-clone/shopee/services/payment/internal/infrastructure/redis"
	"github.com/shopee-clone/shopee/services/payment/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type PaymentService struct {
	cfg           *config.Config
	paymentRepo   *mysql.PaymentRepository
	redisStore    *redisinfra.Store
	kafkaProducer *kafka.Producer
}

func NewPaymentService(cfg *config.Config, paymentRepo *mysql.PaymentRepository, redisStore *redisinfra.Store, kafkaProducer *kafka.Producer) *PaymentService {
	return &PaymentService{cfg: cfg, paymentRepo: paymentRepo, redisStore: redisStore, kafkaProducer: kafkaProducer}
}

type AuthorizePaymentRequest struct {
	OrderID       string              `json:"order_id" validate:"required"`
	UserID        string              `json:"user_id" validate:"required"`
	Amount        int64               `json:"amount" validate:"required"`
	Currency      string              `json:"currency"`
	PaymentMethod domain.PaymentMethod `json:"payment_method" validate:"required"`
	IdempotencyKey string             `json:"idempotency_key" validate:"required"`
	Metadata      json.RawMessage     `json:"metadata,omitempty"`
}

func (s *PaymentService) AuthorizePayment(ctx context.Context, req *AuthorizePaymentRequest) (*domain.Payment, error) {
	ctx, span := otel.Tracer("shopee-payment").Start(ctx, "PaymentService.AuthorizePayment")
	defer span.End()

	start := time.Now()
	defer func() { metrics.PaymentAuthorizationLatency.WithLabelValues(s.cfg.Payment.DefaultPSP).Observe(time.Since(start).Seconds()) }()

	// Idempotency check
	if req.IdempotencyKey != "" {
		existingID, err := s.redisStore.CheckIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && existingID != "" {
			metrics.DuplicatePreventionCount.Inc()
			return s.paymentRepo.FindByID(ctx, existingID)
		}
		existing, err := s.paymentRepo.FindByIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && existing != nil {
			metrics.DuplicatePreventionCount.Inc()
			return existing, nil
		}
	}

	// [SECURITY] Anti-double-charge: Use distributed lock BEFORE checking existing payment.
	// This prevents race conditions where two concurrent requests both pass the check.
	locked, err := s.redisStore.AcquirePaymentLock(ctx, req.OrderID, 30*time.Second)
	if err != nil || !locked {
		return nil, fmt.Errorf("failed to acquire payment lock")
	}
	defer s.redisStore.ReleasePaymentLock(ctx, req.OrderID)

	// Check if payment already exists for this order (inside lock)
	existingPayment, err := s.paymentRepo.FindByOrderID(ctx, req.OrderID)
	if err == nil && existingPayment != nil && !existingPayment.IsTerminal() {
		span.SetStatus(codes.Error, "double charge detected")
		return nil, domain.ErrDoubleChargeDetected
	}

	currency := req.Currency
	if currency == "" { currency = "SGD" }

	payment := domain.NewPayment(req.OrderID, req.UserID, req.Amount, currency, req.PaymentMethod, s.cfg.Payment.DefaultPSP, req.IdempotencyKey)
	payment.Metadata = req.Metadata

	// Fraud check (async hook)
	fraudResult := domain.NewFraudCheckResult(payment.ID, req.UserID, 10, false)
	s.paymentRepo.SaveFraudCheck(ctx, fraudResult)

	// Simulate PSP authorization
	payment.PSPTransactionID = fmt.Sprintf("psp-tx-%s", payment.ID[:8])
	if err := payment.TransitionTo(domain.PaymentStatusAuthorized); err != nil {
		return nil, err
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Store idempotency
	if req.IdempotencyKey != "" {
		rec := domain.NewIdempotencyRecord(payment.ID, s.cfg.Payment.IdempotencyTTL)
		rec.Key = req.IdempotencyKey
		s.paymentRepo.SaveIdempotencyKey(ctx, rec)
		s.redisStore.StoreIdempotencyKey(ctx, req.IdempotencyKey, payment.ID, s.cfg.Payment.IdempotencyTTL)
	}

	// Publish event
	event := domain.NewPaymentEvent(payment, domain.EventPaymentAuthorized, req.Metadata)
	payload, _ := json.Marshal(event)
	s.paymentRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("payment", payment.ID, string(domain.EventPaymentAuthorized), payload))
	if s.kafkaProducer != nil {
		s.kafkaProducer.PublishEvent(ctx, event)
	}

	metrics.PaymentsAuthorizedTotal.WithLabelValues(s.cfg.Payment.DefaultPSP, string(req.PaymentMethod)).Inc()
	metrics.ActivePayments.WithLabelValues(string(domain.PaymentStatusAuthorized)).Inc()

	span.SetAttributes(attribute.String("payment_id", payment.ID), attribute.Int64("amount", payment.Amount))
	zap.L().Info("payment authorized", zap.String("payment_id", payment.ID), zap.String("order_id", req.OrderID))
	return payment, nil
}

func (s *PaymentService) CapturePayment(ctx context.Context, paymentID, actorID string) (*domain.Payment, error) {
	ctx, span := otel.Tracer("shopee-payment").Start(ctx, "PaymentService.CapturePayment")
	defer span.End()

	start := time.Now()
	defer func() { metrics.PaymentCaptureLatency.WithLabelValues(s.cfg.Payment.DefaultPSP).Observe(time.Since(start).Seconds()) }()

	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil { return nil, err }

	if err := payment.TransitionTo(domain.PaymentStatusCaptured); err != nil {
		return nil, err
	}

	if err := s.paymentRepo.UpdateStatus(ctx, paymentID, domain.PaymentStatusCaptured, payment.Version-1); err != nil {
		return nil, err
	}

	event := domain.NewPaymentEvent(payment, domain.EventPaymentCaptured, nil)
	payload, _ := json.Marshal(event)
	s.paymentRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("payment", payment.ID, string(domain.EventPaymentCaptured), payload))
	if s.kafkaProducer != nil { s.kafkaProducer.PublishEvent(ctx, event) }

	metrics.PaymentsCapturedTotal.WithLabelValues(s.cfg.Payment.DefaultPSP).Inc()
	return payment, nil
}

func (s *PaymentService) RefundPayment(ctx context.Context, paymentID, reason, idempotencyKey string, amount int64) (*domain.Refund, error) {
	ctx, span := otel.Tracer("shopee-payment").Start(ctx, "PaymentService.RefundPayment")
	defer span.End()

	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil { return nil, err }

	if payment.Status != domain.PaymentStatusCaptured && payment.Status != domain.PaymentStatusPartialRefund {
		return nil, domain.ErrRefundNotAllowed
	}

	// [FIX BUG-9] Validate amount is positive and doesn't exceed remaining
	if amount <= 0 {
		return nil, fmt.Errorf("refund amount must be positive, got %d", amount)
	}
	if amount > payment.RemainingAmount() {
		return nil, domain.ErrRefundAmountExceeded
	}

	refund := domain.NewRefund(paymentID, payment.OrderID, payment.Currency, reason, idempotencyKey, amount)
	if err := s.paymentRepo.SaveRefund(ctx, refund); err != nil {
		return nil, err
	}

	payment.AmountRefunded += amount
	newStatus := domain.PaymentStatusPartialRefund
	if payment.AmountRefunded >= payment.Amount {
		newStatus = domain.PaymentStatusRefunded
	}
	payment.TransitionTo(newStatus)
	s.paymentRepo.Update(ctx, payment)

	event := domain.NewPaymentEvent(payment, domain.EventPaymentRefunded, nil)
	payload, _ := json.Marshal(event)
	s.paymentRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("payment", payment.ID, string(domain.EventPaymentRefunded), payload))
	if s.kafkaProducer != nil { s.kafkaProducer.PublishEvent(ctx, event) }

	metrics.RefundsProcessed.WithLabelValues("success").Inc()
	return refund, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	return s.paymentRepo.FindByID(ctx, paymentID)
}

func (s *PaymentService) HandleWebhook(ctx context.Context, pspProvider, eventType string, payload []byte, signature, idempotencyKey string) error {
	ctx, span := otel.Tracer("shopee-payment").Start(ctx, "PaymentService.HandleWebhook")
	defer span.End()

	start := time.Now()
	defer func() { metrics.WebhookLatency.WithLabelValues(pspProvider, eventType).Observe(time.Since(start).Seconds()) }()

	// Replay protection
	isReplay, err := s.redisStore.CheckWebhookReplay(ctx, idempotencyKey)
	if err == nil && isReplay {
		metrics.ReplayAttackCount.Inc()
		return domain.ErrWebhookReplayDetected
	}

	// Verify signature
	if !domain.VerifyWebhookSignature(payload, signature, s.cfg.Payment.WebhookSecret) {
		return domain.ErrInvalidWebhookSignature
	}

	// Store webhook event
	webhookEvent := domain.NewWebhookEvent(pspProvider, eventType, payload, signature, idempotencyKey)
	if err := s.paymentRepo.SaveWebhookEvent(ctx, webhookEvent); err != nil {
		return err
	}
	s.redisStore.MarkWebhookProcessed(ctx, idempotencyKey, 24*time.Hour)

	// Process based on event type
	switch eventType {
	case "payment.authorized":
		// Update payment status
	case "payment.captured":
		// Update payment status
	case "payment.failed":
		// Mark payment as failed
	}

	metrics.WebhookProcessed.WithLabelValues(pspProvider, eventType).Inc()
	return nil
}

func (s *PaymentService) ProcessOutboxEvents(ctx context.Context) error {
	events, err := s.paymentRepo.GetUnprocessedOutboxEvents(ctx, 100)
	if err != nil { return err }
	for _, event := range events {
		var paymentEvent domain.PaymentEvent
		if err := json.Unmarshal(event.Payload, &paymentEvent); err != nil { continue }
		if err := s.kafkaProducer.PublishEvent(ctx, &paymentEvent); err != nil { continue }
		s.paymentRepo.MarkOutboxEventProcessed(ctx, event.ID)
	}
	return nil
}
