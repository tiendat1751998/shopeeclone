package events

import (
	"context"
	"time"
)

type EventType string

const (
	ShipmentCreated       EventType = "shipment.created"
	ShipmentStatusChanged EventType = "shipment.status_changed"
	TrackingUpdated       EventType = "tracking.updated"
	CourierWebhookReceived EventType = "courier.webhook_received"
	DispatchCreated       EventType = "dispatch.created"
	DispatchAssigned      EventType = "dispatch.assigned"
	DispatchCompleted     EventType = "dispatch.completed"
	PickupCompleted       EventType = "pickup.completed"
	PickupFailed          EventType = "pickup.failed"
	FulfillmentPacked     EventType = "fulfillment.packed"
	FulfillmentShipped    EventType = "fulfillment.shipped"
	EstimationCalculated  EventType = "estimation.calculated"
)

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Source    string    `json:"source"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

type Producer interface {
	Publish(ctx context.Context, event Event) error
	PublishBatch(ctx context.Context, events []Event) error
	Close() error
}

type Consumer interface {
	Consume(ctx context.Context, handler func(Event) error) error
	Close() error
}
