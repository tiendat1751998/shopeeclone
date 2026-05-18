package application

import (
	"context"
	"time"

	"github.com/shopee-clone/shopee/services/order/internal/domain"
	"github.com/shopee-clone/shopee/services/order/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func (s *OrderService) TriggerReconciliation(ctx context.Context, orderID string, rtype domain.ReconciliationType) (*domain.OrderReconciliation, error) {
	ctx, span := otel.Tracer("shopee-order").Start(ctx, "OrderService.TriggerReconciliation")
	defer span.End()

	rec := domain.NewOrderReconciliation(orderID, rtype)
	if err := s.orderRepo.SaveReconciliation(ctx, rec); err != nil {
		return nil, err
	}

	// Publish reconciliation event
	event := &domain.OrderEvent{
		OrderID:   orderID,
		EventType: domain.EventOrderReconciliationTriggered,
		Timestamp: time.Now().UTC(),
	}
	if s.kafkaProducer != nil {
		s.kafkaProducer.PublishEvent(ctx, event)
	}

	return rec, nil
}

func (s *OrderService) RunReconciliation(ctx context.Context) error {
	ctx, span := otel.Tracer("shopee-order").Start(ctx, "OrderService.RunReconciliation")
	defer span.End()

	recs, err := s.orderRepo.GetPendingReconciliations(ctx, 100)
	if err != nil {
		return err
	}

	for _, rec := range recs {
		start := time.Now()
		s.reconcileOrder(ctx, rec)
		metrics.ReconciliationLatency.WithLabelValues(string(rec.ReconciliationType)).Observe(time.Since(start).Seconds())
	}

	return nil
}

func (s *OrderService) reconcileOrder(ctx context.Context, rec *domain.OrderReconciliation) {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("panic in reconciliation", zap.Any("recover", r))
		}
	}()

	rec.IncrementRetry()
	rec.MarkChecked()

	order, err := s.orderRepo.FindByID(ctx, rec.OrderID)
	if err != nil {
		s.orderRepo.UpdateReconciliationStatus(ctx, rec.ID, domain.ReconciliationStatusFailed, rec.RetryCount)
		metrics.ReconciliationFailures.WithLabelValues(string(rec.ReconciliationType)).Inc()
		return
	}

	var status domain.ReconciliationStatus

	switch rec.ReconciliationType {
	case domain.ReconciliationTypePayment:
		status = s.reconcilePayment(ctx, order)
	case domain.ReconciliationTypeInventory:
		status = s.reconcileInventory(ctx, order)
	case domain.ReconciliationTypeShipment:
		status = s.reconcileShipment(ctx, order)
	default:
		status = domain.ReconciliationStatusFailed
	}

	s.orderRepo.UpdateReconciliationStatus(ctx, rec.ID, status, rec.RetryCount)

	if status == domain.ReconciliationStatusFailed && rec.CanRetry() {
		metrics.ReconciliationFailures.WithLabelValues(string(rec.ReconciliationType)).Inc()
	}
}

func (s *OrderService) reconcilePayment(ctx context.Context, order *domain.Order) domain.ReconciliationStatus {
	// In production: call payment service to verify payment status
	// For now, mark as matched if order is paid or beyond
	if order.Status == domain.OrderStatusPaid ||
		order.Status == domain.OrderStatusProcessing ||
		order.Status == domain.OrderStatusPacked ||
		order.Status == domain.OrderStatusShipped ||
		order.Status == domain.OrderStatusDelivered ||
		order.Status == domain.OrderStatusCompleted {
		return domain.ReconciliationStatusMatched
	}
	if order.Status == domain.OrderStatusCancelled || order.Status == domain.OrderStatusRefunded {
		return domain.ReconciliationStatusMatched
	}
	return domain.ReconciliationStatusPending
}

func (s *OrderService) reconcileInventory(ctx context.Context, order *domain.Order) domain.ReconciliationStatus {
	// In production: call inventory service to verify stock deductions
	// Check if order has reached paid status or beyond (using valid state progression)
	if order.Status == domain.OrderStatusPaid ||
		order.Status == domain.OrderStatusProcessing ||
		order.Status == domain.OrderStatusPacked ||
		order.Status == domain.OrderStatusShipped ||
		order.Status == domain.OrderStatusDelivered ||
		order.Status == domain.OrderStatusCompleted {
		return domain.ReconciliationStatusMatched
	}
	return domain.ReconciliationStatusPending
}

func (s *OrderService) reconcileShipment(ctx context.Context, order *domain.Order) domain.ReconciliationStatus {
	// In production: call shipment service to verify shipment status
	if order.Status == domain.OrderStatusShipped ||
		order.Status == domain.OrderStatusDelivered ||
		order.Status == domain.OrderStatusCompleted {
		return domain.ReconciliationStatusMatched
	}
	if order.Status == domain.OrderStatusCancelled || order.Status == domain.OrderStatusRefunded {
		return domain.ReconciliationStatusMatched
	}
	return domain.ReconciliationStatusPending
}
