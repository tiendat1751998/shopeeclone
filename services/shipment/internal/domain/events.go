package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ShipmentEventType string

const (
	EventShipmentCreated   ShipmentEventType = "shipment.created"
	EventShipmentBooked    ShipmentEventType = "shipment.booked"
	EventShipmentPickedUp  ShipmentEventType = "shipment.picked_up"
	EventShipmentInTransit ShipmentEventType = "shipment.in_transit"
	EventShipmentDelivered ShipmentEventType = "shipment.delivered"
	EventShipmentFailed    ShipmentEventType = "shipment.failed"
	EventShipmentReturned  ShipmentEventType = "shipment.returned"
	EventTrackingUpdated   ShipmentEventType = "shipment.tracking_updated"
)

type ShipmentEvent struct {
	ShipmentID    string           `json:"shipment_id"`
	OrderID       string           `json:"order_id"`
	Status        ShipmentStatus   `json:"status"`
	EventType     ShipmentEventType `json:"event_type"`
	Metadata      json.RawMessage  `json:"metadata,omitempty"`
	Timestamp     time.Time        `json:"timestamp"`
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
