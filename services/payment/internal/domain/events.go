package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PaymentEventType string

const (
	EventPaymentAuthorized PaymentEventType = "payment.authorized"
	EventPaymentCaptured   PaymentEventType = "payment.captured"
	EventPaymentFailed     PaymentEventType = "payment.failed"
	EventPaymentExpired    PaymentEventType = "payment.expired"
	EventPaymentRefunded   PaymentEventType = "payment.refunded"
	EventPaymentReconciled PaymentEventType = "payment.reconciled"
	EventWebhookReceived   PaymentEventType = "payment.webhook_received"
)

type PaymentEvent struct {
	PaymentID   string           `json:"payment_id"`
	OrderID     string           `json:"order_id"`
	UserID      string           `json:"user_id"`
	Amount      int64            `json:"amount"`
	Currency    string           `json:"currency"`
	Status      PaymentStatus    `json:"status"`
	EventType   PaymentEventType `json:"event_type"`
	Metadata    json.RawMessage  `json:"metadata,omitempty"`
	Timestamp   time.Time        `json:"timestamp"`
}

func NewPaymentEvent(payment *Payment, eventType PaymentEventType, metadata json.RawMessage) *PaymentEvent {
	return &PaymentEvent{
		PaymentID: payment.ID,
		OrderID:   payment.OrderID,
		UserID:    payment.UserID,
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		Status:    payment.Status,
		EventType: eventType,
		Metadata:  metadata,
		Timestamp: time.Now().UTC(),
	}
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
		ID:            uuid.New().String(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       payload,
		CreatedAt:     time.Now().UTC(),
		Processed:     false,
	}
}
