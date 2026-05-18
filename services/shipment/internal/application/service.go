package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopee-clone/shopee/services/shipment/internal/config"
	"github.com/shopee-clone/shopee/services/shipment/internal/domain"
	"github.com/shopee-clone/shopee/services/shipment/internal/infrastructure/kafka"
	"github.com/shopee-clone/shopee/services/shipment/internal/infrastructure/mysql"
	redisinfra "github.com/shopee-clone/shopee/services/shipment/internal/infrastructure/redis"
	"github.com/shopee-clone/shopee/services/shipment/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type ShipmentService struct {
	cfg           *config.Config
	shipmentRepo  *mysql.ShipmentRepository
	redisStore    *redisinfra.Store
	kafkaProducer *kafka.Producer
}

func NewShipmentService(cfg *config.Config, repo *mysql.ShipmentRepository, store *redisinfra.Store, producer *kafka.Producer) *ShipmentService {
	return &ShipmentService{cfg: cfg, shipmentRepo: repo, redisStore: store, kafkaProducer: producer}
}

type CreateShipmentRequest struct {
	OrderID        string      `json:"order_id" validate:"required"`
	UserID         string      `json:"user_id" validate:"required"`
	CarrierID      string      `json:"carrier_id"`
	IdempotencyKey string      `json:"idempotency_key" validate:"required"`
	Origin         domain.Address `json:"origin" validate:"required"`
	Destination    domain.Address `json:"destination" validate:"required"`
	Weight         float64     `json:"weight"`
	Currency       string      `json:"currency"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
}

func (s *ShipmentService) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*domain.Shipment, error) {
	ctx, span := otel.Tracer("shopee-shipment").Start(ctx, "ShipmentService.CreateShipment")
	defer span.End()

	// Idempotency check
	if req.IdempotencyKey != "" {
		existingID, err := s.redisStore.CheckIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && existingID != "" {
			return s.shipmentRepo.FindByID(ctx, existingID)
		}
		existing, err := s.shipmentRepo.FindByIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && existing != nil { return existing, nil }
	}

	carrierID := req.CarrierID
	if carrierID == "" { carrierID = s.cfg.Shipment.DefaultCarrier }
	currency := req.Currency
	if currency == "" { currency = "SGD" }

	shipment := domain.NewShipment(req.OrderID, req.UserID, carrierID, req.IdempotencyKey, currency, req.Origin, req.Destination, req.Weight)
	shipment.Metadata = req.Metadata

	// Acquire lock
	locked, err := s.redisStore.AcquireShipmentLock(ctx, req.OrderID, 30*time.Second)
	if err != nil || !locked { return nil, fmt.Errorf("failed to acquire shipment lock") }
	defer s.redisStore.ReleaseShipmentLock(ctx, req.OrderID)

	if err := s.shipmentRepo.Create(ctx, shipment); err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	if req.IdempotencyKey != "" {
		rec := domain.NewIdempotencyRecord(shipment.ID, s.cfg.Shipment.IdempotencyTTL)
		rec.Key = req.IdempotencyKey
		s.shipmentRepo.SaveIdempotencyKey(ctx, rec)
		s.redisStore.StoreIdempotencyKey(ctx, req.IdempotencyKey, shipment.ID, s.cfg.Shipment.IdempotencyTTL)
	}

	event := &domain.ShipmentEvent{ShipmentID: shipment.ID, OrderID: req.OrderID, Status: shipment.Status, EventType: domain.EventShipmentCreated, Timestamp: time.Now().UTC()}
	payload, _ := json.Marshal(event)
	s.shipmentRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("shipment", shipment.ID, string(domain.EventShipmentCreated), payload))
	if s.kafkaProducer != nil { s.kafkaProducer.PublishEvent(ctx, event) }

	metrics.ShipmentsCreatedTotal.WithLabelValues(carrierID).Inc()
	metrics.ActiveShipments.WithLabelValues(string(domain.ShipmentStatusPending)).Inc()

	span.SetAttributes(attribute.String("shipment_id", shipment.ID))
	zap.L().Info("shipment created", zap.String("shipment_id", shipment.ID), zap.String("order_id", req.OrderID))
	return shipment, nil
}

func (s *ShipmentService) UpdateStatus(ctx context.Context, shipmentID string, target domain.ShipmentStatus, actorID, reason string) (*domain.Shipment, error) {
	ctx, span := otel.Tracer("shopee-shipment").Start(ctx, "ShipmentService.UpdateStatus")
	defer span.End()

	shipment, err := s.shipmentRepo.FindByID(ctx, shipmentID)
	if err != nil { return nil, err }

	if err := shipment.TransitionTo(target); err != nil { return nil, err }

	if err := s.shipmentRepo.UpdateStatus(ctx, shipmentID, target, shipment.Version-1); err != nil {
		return nil, err
	}

	event := &domain.ShipmentEvent{ShipmentID: shipment.ID, OrderID: shipment.OrderID, Status: target, EventType: domain.ShipmentEventType("shipment."+string(target)), Timestamp: time.Now().UTC()}
	payload, _ := json.Marshal(event)
	s.shipmentRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("shipment", shipment.ID, string(event.EventType), payload))
	if s.kafkaProducer != nil { s.kafkaProducer.PublishEvent(ctx, event) }

	metrics.ShipmentTransitionLatency.WithLabelValues(string(shipment.Status), string(target)).Observe(0)
	return shipment, nil
}

func (s *ShipmentService) GetShipment(ctx context.Context, shipmentID string) (*domain.Shipment, error) {
	return s.shipmentRepo.FindByID(ctx, shipmentID)
}

func (s *ShipmentService) GetTrackingHistory(ctx context.Context, shipmentID string) ([]*domain.TrackingEvent, error) {
	return s.shipmentRepo.GetTrackingHistory(ctx, shipmentID)
}

func (s *ShipmentService) HandleWebhook(ctx context.Context, provider, eventType string, payload []byte, signature, idempotencyKey string) error {
	isReplay, err := s.redisStore.CheckWebhookReplay(ctx, idempotencyKey)
	if err == nil && isReplay { metrics.WebhookReplayCount.Inc(); return domain.ErrWebhookReplayDetected }

	s.shipmentRepo.SaveWebhookEvent(ctx, provider, eventType, payload, signature, idempotencyKey)
	s.redisStore.MarkWebhookProcessed(ctx, idempotencyKey, 24*time.Hour)
	return nil
}

func (s *ShipmentService) ProcessOutboxEvents(ctx context.Context) error {
	events, err := s.shipmentRepo.GetUnprocessedOutboxEvents(ctx, 100)
	if err != nil { return err }
	for _, event := range events {
		var shipmentEvent domain.ShipmentEvent
		if err := json.Unmarshal(event.Payload, &shipmentEvent); err != nil { continue }
		if err := s.kafkaProducer.PublishEvent(ctx, &shipmentEvent); err != nil { continue }
		s.shipmentRepo.MarkOutboxEventProcessed(ctx, event.ID)
	}
	return nil
}
