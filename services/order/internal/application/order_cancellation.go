package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopee-clone/shopee/services/order/internal/domain"
	"github.com/shopee-clone/shopee/services/order/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type CancelOrderRequest struct {
	OrderID    string                  `json:"order_id" validate:"required"`
	Reason     string                  `json:"reason" validate:"required"`
	CancelledBy string                 `json:"cancelled_by" validate:"required"`
	CancelledType domain.CancellationType `json:"cancelled_type" validate:"required"`
}

func (s *OrderService) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*domain.Order, error) {
	ctx, span := otel.Tracer("shopee-order").Start(ctx, "OrderService.CancelOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("order_id", req.OrderID),
		attribute.String("cancelled_by", req.CancelledBy),
		attribute.String("cancellation_type", string(req.CancelledType)),
	)

	start := time.Now()
	defer func() {
		metrics.OrderCancellationLatency.Observe(time.Since(start).Seconds())
	}()

	// Acquire distributed lock
	locked, err := s.redisStore.AcquireTransitionLock(ctx, req.OrderID, 10*time.Second)
	if err != nil || !locked {
		return nil, fmt.Errorf("failed to acquire cancellation lock: %w", err)
	}
	defer s.redisStore.ReleaseTransitionLock(ctx, req.OrderID)

	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	// Validate cancellation
	if !order.IsCancellable() {
		span.SetStatus(codes.Error, "order not cancellable")
		return nil, domain.ErrOrderNotCancellable
	}

	// Perform cancellation transition
	lifecycleEvent, err := order.TransitionTo(domain.OrderStatusCancelled, req.CancelledBy, string(req.CancelledType), req.Reason)
	if err != nil {
		return nil, err
	}

	// Persist
	// TransitionTo already incremented Version, so we use the current Version as expected
	if err := s.orderRepo.UpdateStatus(ctx, req.OrderID, domain.OrderStatusCancelled, order.Version); err != nil {
		return nil, err
	}

	// Save lifecycle event
	s.orderRepo.SaveLifecycleEvent(ctx, lifecycleEvent)

	// Record cancellation
	cancellation := domain.NewOrderCancellation(
		req.OrderID, req.Reason, req.CancelledBy, req.CancelledType, order.TotalAmount,
	)
	s.orderRepo.SaveCancellation(ctx, cancellation)

	// Invalidate cache
	s.redisStore.InvalidateOrderCache(ctx, req.OrderID)

	// Publish event
	event := domain.NewOrderEvent(order, domain.EventOrderCancelled, nil)
	eventPayload, _ := json.Marshal(event)
	outboxEvent := domain.NewOutboxEvent("order", order.ID, string(domain.EventOrderCancelled), eventPayload)
	s.orderRepo.SaveOutboxEvent(ctx, outboxEvent)

	if s.kafkaProducer != nil {
		s.kafkaProducer.PublishEvent(ctx, event)
	}

	// Trigger compensation workflows asynchronously
	compCtx, compCancel := context.WithTimeout(context.Background(), 30*time.Second)
	go func() {
		defer compCancel()
		s.triggerCompensation(compCtx, order, cancellation)
	}()

	// Update metrics
	metrics.OrdersCancelledTotal.WithLabelValues(string(req.CancelledType)).Inc()
	metrics.ActiveOrdersByStatus.WithLabelValues(string(order.Status)).Dec()

	zap.L().Info("order cancelled",
		zap.String("order_id", req.OrderID),
		zap.String("reason", req.Reason),
		zap.String("cancelled_by", req.CancelledBy),
	)

	return order, nil
}

func (s *OrderService) triggerCompensation(ctx context.Context, order *domain.Order, cancellation *domain.OrderCancellation) {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("panic in compensation", zap.Any("recover", r), zap.String("order_id", order.ID))
		}
	}()

	// Compensation: release inventory reservations, trigger refund if paid
	zap.L().Info("triggering compensation for cancelled order",
		zap.String("order_id", order.ID),
		zap.String("cancellation_id", cancellation.ID),
	)

	// Update compensation status
	s.orderRepo.UpdateCancellationCompensation(ctx, cancellation.ID, domain.CompensationInProgress)

	// In a real system, this would:
	// 1. Call inventory service to release reservations
	// 2. Call payment service to trigger refund if order was paid
	// 3. Call shipment service to cancel shipment if order was shipped

	s.orderRepo.UpdateCancellationCompensation(ctx, cancellation.ID, domain.CompensationCompleted)
}
