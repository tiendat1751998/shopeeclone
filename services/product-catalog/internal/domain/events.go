package domain

import "time"

type CatalogEvent struct {
	ID            string      `json:"id"`
	EventType     string      `json:"event_type"`
	AggregateType string      `json:"aggregate_type"`
	AggregateID   string      `json:"aggregate_id"`
	Payload       interface{} `json:"payload"`
	CreatedAt     time.Time   `json:"created_at"`
}

// OutboxEvent represents an event persisted in the outbox table for reliable delivery.
type OutboxEvent struct {
	ID            string    `db:"event_id" json:"event_id"`
	AggregateType string    `db:"aggregate_type" json:"aggregate_type"`
	AggregateID   string    `db:"aggregate_id" json:"aggregate_id"`
	EventType     string    `db:"event_type" json:"event_type"`
	Payload       string    `db:"payload" json:"payload"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Processed     bool      `db:"processed" json:"processed"`
}

const (
	EventProductCreated  = "product.created"
	EventProductUpdated  = "product.updated"
	EventProductArchived = "product.archived"
	EventSKUUpdated      = "sku.updated"
	EventCategoryUpdated = "category.updated"
	EventMediaUpdated    = "media.updated"
	EventIndexingTriggered = "indexing.triggered"
)
