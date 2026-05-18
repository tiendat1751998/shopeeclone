package domain

import "time"

// CartEvent represents events emitted by the cart service
type CartEvent struct {
	ID            string      `json:"id"`
	EventType     string      `json:"event_type"`
	AggregateType string      `json:"aggregate_type"`
	AggregateID   string      `json:"aggregate_id"`
	Payload       interface{} `json:"payload"`
	CreatedAt     time.Time   `json:"created_at"`
}

const (
	EventCartCreated      = "cart.created"
	EventCartUpdated      = "cart.updated"
	EventCartMerged       = "cart.merged"
	EventCartExpired      = "cart.expired"
	EventItemAdded        = "cart.item_added"
	EventItemRemoved      = "cart.item_removed"
	EventItemUpdated      = "cart.item_updated"
	EventCheckoutPrepared = "cart.checkout_prepared"
	EventCartCleared      = "cart.cleared"
)

type CartUpdatedPayload struct {
	CartID    string `json:"cart_id"`
	UserID    string `json:"user_id"`
	ItemCount int    `json:"item_count"`
	Subtotal  int64  `json:"subtotal"`
}

type CartMergedPayload struct {
	SourceCartID string `json:"source_cart_id"`
	TargetCartID string `json:"target_cart_id"`
	UserID       string `json:"user_id"`
	ItemsMerged  int    `json:"items_merged"`
}

type ItemAddedPayload struct {
	CartID      string `json:"cart_id"`
	SKU         string `json:"sku"`
	ProductName string `json:"product_name"`
	ShopID      string `json:"shop_id"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`
}

type CheckoutPreparedPayload struct {
	CartID    string `json:"cart_id"`
	UserID    string `json:"user_id"`
	Subtotal  int64  `json:"subtotal"`
	ItemCount int    `json:"item_count"`
}
