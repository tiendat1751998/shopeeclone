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
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
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
	OrderID        string              `json:"order_id" validate:"required"`
	UserID         string              `json:"user_id" validate:"required"`
	Amount         int64               `json:"amount" validate:"required"`
	Currency       string              `json:"currency"`
	PaymentMethod  domain.PaymentMethod `json:"payment_method" validate:"required"`
	IdempotencyKey string              `json:"idempotency_key" validate:"required"`
	Metadata       json.RawMessage     `json:"metadata,omitempty"`
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

	// [FIX A2] Fraud check - MUST handle error
	fraudResult := domain.NewFraudCheckResult(payment.ID, req.UserID, 10, false)
	if err := s.paymentRepo.SaveFraudCheck(ctx, fraudResult); err != nil {
		observability.LogWithTrace(ctx).Error("failed to save fraud check", zap.Error(err))
		// Don't fail the payment, but log for investigation
	}

	// Simulate PSP authorization
	payment.PSPTransactionID = fmt.Sprintf("psp-tx-%s", payment.ID[:8])
	if err := payment.TransitionTo(domain.PaymentStatusAuthorized); err != nil {
		return nil, err
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// [FIX A2] Store idempotency - MUST handle error
	if req.IdempotencyKey != "" {
		rec := domain.NewIdempotencyRecord(payment.ID, s.cfg.Payment.IdempotencyTTL)
		rec.Key = req.IdempotencyKey
		if err := s.paymentRepo.SaveIdempotencyKey(ctx, rec); err != nil {
			observability.LogWithTrace(ctx).Error("failed to save idempotency key", zap.Error(err))
		}
		if err := s.redisStore.StoreIdempotencyKey(ctx, req.IdempotencyKey, payment.ID, s.cfg.Payment.IdempotencyTTL); err != nil {
			observability.LogWithTrace(ctx).Error("failed to store idempotency key in Redis", zap.Error(err))
		}
	}

	// [FIX A2] Publish event - MUST handle error
	event := domain.NewPaymentEvent(payment, domain.EventPaymentAuthorized, req.Metadata)
	payload, err := json.Marshal(event)
	if err != nil {
		observability.LogWithTrace(ctx).Error("failed to marshal payment event", zap.Error(err))
	} else {
		if err := s.paymentRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("payment", payment.ID, string(domain.EventPaymentAuthorized), payload)); err != nil {
			observability.LogWithTrace(ctx).Error("failed to save outbox event", zap.Error(err))
		}
	}
	if s.kafkaProducer != nil {
		if err := s.kafkaProducer.PublishEvent(ctx, event); err != nil {
			observability.LogWithTrace(ctx).Error("failed to publish payment event to Kafka", zap.Error(err))
		}
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
	payload, err := json.Marshal(event)
	if err != nil {
		observability.LogWithTrace(ctx).Error("failed to marshal capture event", zap.Error(err))
	} else {
		s.paymentRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("payment", payment.ID, string(domain.EventPaymentCaptured), payload))
	}
	if s.kafkaProducer != nil {
		if err := s.kafkaProducer.PublishEvent(ctx, event); err != nil {
			observability.LogWithTrace(ctx).Error("failed to publish capture event", zap.Error(err))
		}
	}

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

	// [FIX A4] Validate amount is positive and doesn't exceed remaining
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

	// [FIX A4] MUST check TransitionTo error
	if err := payment.TransitionTo(newStatus); err != nil {
		return nil, fmt.Errorf("failed to transition payment status: %w", err)
	}

	// [FIX A4] MUST check Update error
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	event := domain.NewPaymentEvent(payment, domain.EventPaymentRefunded, nil)
	payload, err := json.Marshal(event)
	if err != nil {
		observability.LogWithTrace(ctx).Error("failed to marshal refund event", zap.Error(err))
	} else {
		s.paymentRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("payment", payment.ID, string(domain.EventPaymentRefunded), payload))
	}
	if s.kafkaProducer != nil {
		if err := s.kafkaProducer.PublishEvent(ctx, event); err != nil {
			observability.LogWithTrace(ctx).Error("failed to publish refund event", zap.Error(err))
		}
	}

	metrics.RefundsProcessed.WithLabelValues("success").Inc()
	return refund, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	return s.paymentRepo.FindByID(ctx, paymentID)
}

// [FIX A3] Webhook handler - now properly processes PSP events
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

	// [FIX A3] Actually process the webhook event
	switch eventType {
	case "payment.authorized":
		var eventData map[string]interface{}
		if err := json.Unmarshal(payload, &eventData); err != nil {
			return fmt.Errorf("failed to unmarshal webhook payload: %w", err)
		}
		paymentID, _ := eventData["payment_id"].(string)
		if paymentID != "" {
			if err := s.markPaymentAuthorized(ctx, paymentID); err != nil {
				return fmt.Errorf("failed to mark payment authorized: %w", err)
			}
		}
	case "payment.captured":
		var eventData map[string]interface{}
		if err := json.Unmarshal(payload, &eventData); err != nil {
			return fmt.Errorf("failed to unmarshal webhook payload: %w", err)
		}
		paymentID, _ := eventData["payment_id"].(string)
		if paymentID != "" {
			if _, err := s.CapturePayment(ctx, paymentID, "psp_webhook"); err != nil {
				return fmt.Errorf("failed to capture payment: %w", err)
			}
		}
	case "payment.failed":
		var eventData map[string]interface{}
		if err := json.Unmarshal(payload, &eventData); err != nil {
			return fmt.Errorf("failed to unmarshal webhook payload: %w", err)
		}
		paymentID, _ := eventData["payment_id"].(string)
		if paymentID != "" {
			if err := s.markPaymentFailed(ctx, paymentID); err != nil {
				return fmt.Errorf("failed to mark payment failed: %w", err)
			}
		}
	default:
		observability.LogWithTrace(ctx).Warn("unknown webhook event type", zap.String("type", eventType))
	}

	metrics.WebhookProcessed.WithLabelValues(pspProvider, eventType).Inc()
	return nil
}

// markPaymentAuthorized updates payment status from PSP webhook
func (s *PaymentService) markPaymentAuthorized(ctx context.Context, paymentID string) error {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil { return err }
	if payment == nil { return fmt.Errorf("payment not found: %s", paymentID) }
	if err := payment.TransitionTo(domain.PaymentStatusAuthorized); err != nil { return err }
	return s.paymentRepo.UpdateStatus(ctx, paymentID, domain.PaymentStatusAuthorized, payment.Version-1)
}

// markPaymentFailed updates payment status from PSP webhook
func (s *PaymentService) markPaymentFailed(ctx context.Context, paymentID string) error {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil { return err }
	if payment == nil { return fmt.Errorf("payment not found: %s", paymentID) }
	if err := payment.TransitionTo(domain.PaymentStatusFailed); err != nil { return err }
	return s.paymentRepo.UpdateStatus(ctx, paymentID, domain.PaymentStatusFailed, payment.Version-1)
}

// [FIX A1] ProcessOutboxEvents - now properly logs errors and tracks failed events
func (s *PaymentService) ProcessOutboxEvents(ctx context.Context) error {
	events, err := s.paymentRepo.GetUnprocessedOutboxEvents(ctx, 100)
	if err != nil { return err }

	for _, event := range events {
		// Mark as processing first (prevents duplicate processing)
		if err := s.paymentRepo.MarkOutboxEventProcessing(ctx, event.ID); err != nil {
			observability.LogWithTrace(ctx).Error("failed to mark outbox event as processing",
				zap.String("event_id", event.ID), zap.Error(err))
			continue
		}

		var paymentEvent domain.PaymentEvent
		if err := json.Unmarshal(event.Payload, &paymentEvent); err != nil {
			observability.LogWithTrace(ctx).Error("failed to unmarshal outbox event payload",
				zap.String("event_id", event.ID), zap.Error(err))
			s.paymentRepo.MarkOutboxEventFailed(ctx, event.ID, err.Error())
			continue
		}

		if err := s.kafkaProducer.PublishEvent(ctx, &paymentEvent); err != nil {
			observability.LogWithTrace(ctx).Error("failed to publish outbox event to Kafka",
				zap.String("event_id", event.ID), zap.Error(err))
			s.paymentRepo.MarkOutboxEventFailed(ctx, event.ID, err.Error())
			continue
		}

		if err := s.paymentRepo.MarkOutboxEventProcessed(ctx, event.ID); err != nil {
			observability.LogWithTrace(ctx).Error("failed to mark outbox event as processed",
				zap.String("event_id", event.ID), zap.Error(err))
		}
	}
	return nil
}
