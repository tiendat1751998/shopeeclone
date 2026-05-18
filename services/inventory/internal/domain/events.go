package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type InventoryEventType string

const (
	EventStockReserved    InventoryEventType = "inventory.stock_reserved"
	EventStockReleased    InventoryEventType = "inventory.stock_released"
	EventStockDeducted    InventoryEventType = "inventory.stock_deducted"
	EventStockReplenished InventoryEventType = "inventory.stock_replenished"
	EventReservationExpired InventoryEventType = "inventory.reservation_expired"
	EventReconciliationTriggered InventoryEventType = "inventory.reconciliation_triggered"
)

type InventoryEvent struct {
	ProductID   string           `json:"product_id"`
	SkuID       string           `json:"sku_id"`
	WarehouseID string           `json:"warehouse_id"`
	Quantity    int              `json:"quantity"`
	EventType   InventoryEventType `json:"event_type"`
	Metadata    json.RawMessage  `json:"metadata,omitempty"`
	Timestamp   time.Time        `json:"timestamp"`
}

type OutboxEvent struct {
	ID            string    `db:"event_id" json:"event_id"`
	AggregateType string    `db:"aggregate_type" json:"aggregate_type"`
	AggregateID   string    `db:"aggregate_id" json:"aggregate_id"`
	EventType     string    `db:"event_type" json:"event_type"`
	Payload       []byte    `db:"payload" json:"payload"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Processed     bool      `db:"processed" json:"processed"`
}

func NewOutboxEvent(aggregateType, aggregateID, eventType string, payload []byte) *OutboxEvent {
	return &OutboxEvent{
		ID: uuid.New().String(), AggregateType: aggregateType, AggregateID: aggregateID,
		EventType: eventType, Payload: payload, CreatedAt: time.Now().UTC(), Processed: false,
	}
}
