package domain

import (
	"encoding/json"
	"time"
)

// EventType identifies the type of domain event.
type EventType string

const (
	EventTypeProductCreated   EventType = "product.created"
	EventTypeProductUpdated   EventType = "product.updated"
	EventTypeProductDeleted   EventType = "product.deleted"
	EventTypeSKUUpdated       EventType = "sku.updated"
	EventTypeCategoryUpdated  EventType = "category.updated"
	EventTypeProductModerated EventType = "product.moderated"
)

// DomainEvent is the base metadata attached to every domain event.
type DomainEvent struct {
	EventID   string    `json:"event_id"`
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

// Marshal serializes the event to JSON
func (e DomainEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// ProductCreatedEvent is emitted when a new product is created.
type ProductCreatedEvent struct {
	DomainEvent
	Product *Product `json:"product"`
}

// ProductUpdatedEvent is emitted when product details change.
type ProductUpdatedEvent struct {
	DomainEvent
	Product        *Product `json:"product"`
	ChangedFields  []string `json:"changed_fields"`
}

// ProductDeletedEvent is emitted when a product is soft-deleted.
type ProductDeletedEvent struct {
	DomainEvent
	Product *Product `json:"product"`
}

// SKUUpdatedEvent is emitted when an SKU's price, stock, or status changes.
type SKUUpdatedEvent struct {
	DomainEvent
	SPUID     string  `json:"spu_id"`
	SKUID     string  `json:"sku_id"`
	Price     float64 `json:"price,omitempty"`
	SalePrice float64 `json:"sale_price,omitempty"`
	Stock     int32   `json:"stock,omitempty"`
	Status    string  `json:"status,omitempty"`
}

// CategoryUpdatedEvent is emitted when a category is modified.
type CategoryUpdatedEvent struct {
	DomainEvent
	Category *Category `json:"category"`
}

// ProductModerationEvent is emitted when a product passes through moderation.
type ProductModerationEvent struct {
	DomainEvent
	SPUID      string           `json:"spu_id"`
	Status     ModerationStatus `json:"status"`
	Reason     string           `json:"reason,omitempty"`
	ReviewerID string           `json:"reviewer_id,omitempty"`
}

// NewProductCreatedEvent creates a ProductCreatedEvent
func NewProductCreatedEvent(product *Product) ProductCreatedEvent {
	return ProductCreatedEvent{
		DomainEvent: DomainEvent{
			Type:      EventTypeProductCreated,
			Timestamp: time.Now().UTC(),
		},
		Product: product,
	}
}

// NewProductUpdatedEvent creates a ProductUpdatedEvent
func NewProductUpdatedEvent(product *Product, changedFields []string) ProductUpdatedEvent {
	return ProductUpdatedEvent{
		DomainEvent: DomainEvent{
			Type:      EventTypeProductUpdated,
			Timestamp: time.Now().UTC(),
		},
		Product:       product,
		ChangedFields: changedFields,
	}
}

// NewProductDeletedEvent creates a ProductDeletedEvent
func NewProductDeletedEvent(product *Product) ProductDeletedEvent {
	return ProductDeletedEvent{
		DomainEvent: DomainEvent{
			Type:      EventTypeProductDeleted,
			Timestamp: time.Now().UTC(),
		},
		Product: product,
	}
}

// NewSKUUpdatedEvent creates an SKUUpdatedEvent
func NewSKUUpdatedEvent(spuID, skuID string, price float64, stock int32, status string) SKUUpdatedEvent {
	return SKUUpdatedEvent{
		DomainEvent: DomainEvent{
			Type:      EventTypeSKUUpdated,
			Timestamp: time.Now().UTC(),
		},
		SPUID:  spuID,
		SKUID:  skuID,
		Price:  price,
		Stock:  stock,
		Status: status,
	}
}

// NewCategoryUpdatedEvent creates a CategoryUpdatedEvent
func NewCategoryUpdatedEvent(category *Category) CategoryUpdatedEvent {
	return CategoryUpdatedEvent{
		DomainEvent: DomainEvent{
			Type:      EventTypeCategoryUpdated,
			Timestamp: time.Now().UTC(),
		},
		Category: category,
	}
}

// NewProductModerationEvent creates a ProductModerationEvent
func NewProductModerationEvent(spuID string, status ModerationStatus, reason, reviewerID string) ProductModerationEvent {
	return ProductModerationEvent{
		DomainEvent: DomainEvent{
			Type:      EventTypeProductModerated,
			Timestamp: time.Now().UTC(),
		},
		SPUID:      spuID,
		Status:     status,
		Reason:     reason,
		ReviewerID: reviewerID,
	}
}

// Marshal methods for each event type
func (e ProductCreatedEvent) Marshal() ([]byte, error)    { return json.Marshal(e) }
func (e ProductUpdatedEvent) Marshal() ([]byte, error)    { return json.Marshal(e) }
func (e ProductDeletedEvent) Marshal() ([]byte, error)    { return json.Marshal(e) }
func (e SKUUpdatedEvent) Marshal() ([]byte, error)        { return json.Marshal(e) }
func (e CategoryUpdatedEvent) Marshal() ([]byte, error)   { return json.Marshal(e) }
func (e ProductModerationEvent) Marshal() ([]byte, error) { return json.Marshal(e) }
