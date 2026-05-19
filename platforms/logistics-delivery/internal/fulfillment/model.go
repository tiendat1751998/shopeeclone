package fulfillment

import "time"

type FulfillmentStatus string

const (
	FulfillmentPending    FulfillmentStatus = "pending"
	FulfillmentProcessing FulfillmentStatus = "processing"
	FulfillmentPacked     FulfillmentStatus = "packed"
	FulfillmentShipped    FulfillmentStatus = "shipped"
	FulfillmentCompleted  FulfillmentStatus = "completed"
	FulfillmentFailed     FulfillmentStatus = "failed"
)

type Fulfillment struct {
	ID          string            `json:"id"`
	ShipmentID  string            `json:"shipment_id"`
	OrderID     string            `json:"order_id"`
	WarehouseID string            `json:"warehouse_id"`
	Status      FulfillmentStatus `json:"status"`
	Items       []FulfillmentItem `json:"items"`
	PackedAt    *time.Time        `json:"packed_at,omitempty"`
	ShippedAt   *time.Time        `json:"shipped_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	Notes       string            `json:"notes,omitempty"`
	ReplayID    string            `json:"replay_id,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type FulfillmentItem struct {
	ID         string `json:"id"`
	ProductID  string `json:"product_id"`
	SKU        string `json:"sku"`
	Quantity   int    `json:"quantity"`
	Location   string `json:"location"`
}
