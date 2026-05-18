package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopee-clone/shopee/services/inventory/internal/config"
	"github.com/shopee-clone/shopee/services/inventory/internal/domain"
	"github.com/shopee-clone/shopee/services/inventory/internal/infrastructure/kafka"
	"github.com/shopee-clone/shopee/services/inventory/internal/infrastructure/mysql"
	redisinfra "github.com/shopee-clone/shopee/services/inventory/internal/infrastructure/redis"
	"github.com/shopee-clone/shopee/services/inventory/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type InventoryService struct {
	cfg           *config.Config
	invRepo       *mysql.InventoryRepository
	redisStore    *redisinfra.Store
	kafkaProducer *kafka.Producer
}

func NewInventoryService(cfg *config.Config, repo *mysql.InventoryRepository, store *redisinfra.Store, producer *kafka.Producer) *InventoryService {
	return &InventoryService{cfg: cfg, invRepo: repo, redisStore: store, kafkaProducer: producer}
}

type ReserveStockRequest struct {
	OrderID        string `json:"order_id" validate:"required"`
	UserID         string `json:"user_id" validate:"required"`
	ProductID      string `json:"product_id" validate:"required"`
	SkuID          string `json:"sku_id" validate:"required"`
	WarehouseID    string `json:"warehouse_id"`
	Quantity       int    `json:"quantity" validate:"required,min=1"`
	IdempotencyKey string `json:"idempotency_key" validate:"required"`
}

func (s *InventoryService) ReserveStock(ctx context.Context, req *ReserveStockRequest) (*domain.Reservation, error) {
	ctx, span := otel.Tracer("shopee-inventory").Start(ctx, "InventoryService.ReserveStock")
	defer span.End()

	start := time.Now()
	defer func() { metrics.ReservationLatency.Observe(time.Since(start).Seconds()) }()

	// Idempotency check
	if req.IdempotencyKey != "" {
		existingID, err := s.redisStore.CheckIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && existingID != "" {
			return s.invRepo.GetReservation(ctx, existingID)
		}
		existing, err := s.invRepo.FindByIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil && existing != nil { return existing, nil }
	}

	// Acquire distributed lock for this SKU
	locked, err := s.redisStore.AcquireStockLock(ctx, req.SkuID, 10*time.Second)
	if err != nil || !locked {
		metrics.ReservationFailures.WithLabelValues("lock_failed").Inc()
		return nil, fmt.Errorf("failed to acquire stock lock")
	}
	defer s.redisStore.ReleaseStockLock(ctx, req.SkuID)

	// Get stock
	stock, err := s.invRepo.GetStock(ctx, req.SkuID, req.WarehouseID)
	if err != nil {
		metrics.ReservationFailures.WithLabelValues("stock_not_found").Inc()
		return nil, err
	}

	// Check available quantity (oversell prevention)
	if stock.AvailableQty < req.Quantity {
		metrics.OversellPreventionCount.Inc()
		metrics.ReservationFailures.WithLabelValues("insufficient_stock").Inc()
		return nil, domain.ErrInsufficientStock
	}

	// Reserve stock
	if err := stock.Reserve(req.Quantity); err != nil {
		return nil, err
	}
	stock.Version++
	if err := s.invRepo.UpdateStock(ctx, stock); err != nil {
		return nil, err
	}

	// Create reservation
	reservation := domain.NewReservation(req.OrderID, req.UserID, req.ProductID, req.SkuID, req.WarehouseID, req.Quantity, s.cfg.Inventory.ReservationTTL, req.IdempotencyKey)
	if err := s.invRepo.SaveReservation(ctx, reservation); err != nil {
		return nil, err
	}

	// Store idempotency
	if req.IdempotencyKey != "" {
		rec := domain.NewIdempotencyRecord(reservation.ID, s.cfg.Inventory.IdempotencyTTL)
		rec.Key = req.IdempotencyKey
		s.invRepo.SaveIdempotencyKey(ctx, rec)
		s.redisStore.StoreIdempotencyKey(ctx, req.IdempotencyKey, reservation.ID, s.cfg.Inventory.IdempotencyTTL)
	}

	// Update cache
	s.redisStore.CacheStock(ctx, req.SkuID, stock.Quantity, stock.ReservedQty, 5*time.Minute)

	// Publish event
	event := &domain.InventoryEvent{ProductID: req.ProductID, SkuID: req.SkuID, WarehouseID: req.WarehouseID, Quantity: req.Quantity, EventType: domain.EventStockReserved, Timestamp: time.Now().UTC()}
	payload, _ := json.Marshal(event)
	s.invRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("inventory", reservation.ID, string(domain.EventStockReserved), payload))
	if s.kafkaProducer != nil { s.kafkaProducer.PublishEvent(ctx, event) }

	span.SetAttributes(attribute.String("reservation_id", reservation.ID), attribute.String("sku_id", req.SkuID), attribute.Int("quantity", req.Quantity))
	zap.L().Info("stock reserved", zap.String("reservation_id", reservation.ID), zap.String("sku_id", req.SkuID), zap.Int("qty", req.Quantity))
	return reservation, nil
}

func (s *InventoryService) ReleaseStock(ctx context.Context, reservationID string) error {
	ctx, span := otel.Tracer("shopee-inventory").Start(ctx, "InventoryService.ReleaseStock")
	defer span.End()

	reservation, err := s.invRepo.GetReservation(ctx, reservationID)
	if err != nil { return err }

	if err := reservation.Release(); err != nil { return err }
	s.invRepo.UpdateReservationStatus(ctx, reservationID, domain.ReservationStatusReleased)

	stock, err := s.invRepo.GetStock(ctx, reservation.SkuID, reservation.WarehouseID)
	if err != nil { return err }
	stock.Release(reservation.Quantity)
	stock.Version++
	s.invRepo.UpdateStock(ctx, stock)

	s.redisStore.CacheStock(ctx, reservation.SkuID, stock.Quantity, stock.ReservedQty, 5*time.Minute)

	event := &domain.InventoryEvent{ProductID: reservation.ProductID, SkuID: reservation.SkuID, WarehouseID: reservation.WarehouseID, Quantity: reservation.Quantity, EventType: domain.EventStockReleased, Timestamp: time.Now().UTC()}
	payload, _ := json.Marshal(event)
	s.invRepo.SaveOutboxEvent(ctx, domain.NewOutboxEvent("inventory", reservation.ID, string(domain.EventStockReleased), payload))
	if s.kafkaProducer != nil { s.kafkaProducer.PublishEvent(ctx, event) }

	return nil
}

func (s *InventoryService) GetStock(ctx context.Context, skuID, warehouseID string) (*domain.Stock, error) {
	return s.invRepo.GetStock(ctx, skuID, warehouseID)
}

func (s *InventoryService) ExpireReservations(ctx context.Context) error {
	reservations, err := s.invRepo.GetExpiredReservations(ctx, 100)
	if err != nil { return err }

	g, ctx := errgroup.WithContext(ctx)
	for _, res := range reservations {
		res := res
		g.Go(func() error {
			if err := s.ReleaseStock(ctx, res.ID); err != nil {
				zap.L().Warn("failed to expire reservation", zap.String("id", res.ID), zap.Error(err))
			}
			return nil
		})
	}
	return g.Wait()
}

func (s *InventoryService) ProcessOutboxEvents(ctx context.Context) error {
	events, err := s.invRepo.GetUnprocessedOutboxEvents(ctx, 100)
	if err != nil { return err }
	for _, event := range events {
		var invEvent domain.InventoryEvent
		if err := json.Unmarshal(event.Payload, &invEvent); err != nil { continue }
		if err := s.kafkaProducer.PublishEvent(ctx, &invEvent); err != nil { continue }
		s.invRepo.MarkOutboxEventProcessed(ctx, event.ID)
	}
	return nil
}
