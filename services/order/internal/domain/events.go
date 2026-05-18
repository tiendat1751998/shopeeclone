package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventOrderCreated              EventType = "order.created"
	EventOrderPaid                 EventType = "order.paid"
	EventOrderProcessing           EventType = "order.processing"
	EventOrderPacked               EventType = "order.packed"
	EventOrderShipped              EventType = "order.shipped"
	EventOrderDelivered            EventType = "order.delivered"
	EventOrderCompleted            EventType = "order.completed"
	EventOrderCancelled            EventType = "order.cancelled"
	EventOrderRefunded             EventType = "order.refunded"
	EventOrderReconciliationTriggered EventType = "order.reconciliation_triggered"
)

type OutboxEvent struct {
	ID             string    `db:"event_id" json:"event_id"`
	AggregateType  string    `db:"aggregate_type" json:"aggregate_type"`
	AggregateID    string    `db:"aggregate_id" json:"aggregate_id"`
	EventType      string    `db:"event_type" json:"event_type"`
	Payload        []byte    `db:"payload" json:"payload"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	Processed      bool      `db:"processed" json:"processed"`
}

type OrderEvent struct {
	OrderID    string          `json:"order_id"`
	OrderNumber string         `json:"order_number"`
	UserID     string          `json:"user_id"`
	SellerID   string          `json:"seller_id"`
	Status     OrderStatus     `json:"status"`
	EventType  EventType       `json:"event_type"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	Timestamp  time.Time       `json:"timestamp"`
}

func NewOrderEvent(order *Order, eventType EventType, metadata json.RawMessage) *OrderEvent {
	return &OrderEvent{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		UserID:      order.UserID,
		SellerID:    order.SellerID,
		Status:      order.Status,
		EventType:   eventType,
		Metadata:    metadata,
		Timestamp:   time.Now().UTC(),
	}
}

func NewOutboxEvent(aggregateType, aggregateID, eventType string, payload []byte) *OutboxEvent {
	return &OutboxEvent{
		ID:            uuid.New().String(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       payload,
		CreatedAt:     time.Now().UTC(),
		Processed:     false,
	}
}
