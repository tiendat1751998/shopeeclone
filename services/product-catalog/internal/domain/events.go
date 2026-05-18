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

const (
	EventProductCreated  = "product.created"
	EventProductUpdated  = "product.updated"
	EventProductArchived = "product.archived"
	EventSKUUpdated      = "sku.updated"
	EventCategoryUpdated = "category.updated"
	EventMediaUpdated    = "media.updated"
	EventIndexingTriggered = "indexing.triggered"
)
