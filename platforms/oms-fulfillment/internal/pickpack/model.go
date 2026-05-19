package pickpack

import "time"

type PickListStatus string

const (
	PickListPending    PickListStatus = "pending"
	PickListInProgress PickListStatus = "in_progress"
	PickListCompleted  PickListStatus = "completed"
)

type PackingStatus string

const (
	PackingPending  PackingStatus = "pending"
	PackingComplete PackingStatus = "complete"
)

type ShipmentStatus string

const (
	ShipmentPending   ShipmentStatus = "pending"
	ShipmentShipped   ShipmentStatus = "shipped"
	ShipmentDelivered ShipmentStatus = "delivered"
)

type PickItem struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	Location  string `json:"location"`
}

type PickList struct {
	ID          string         `json:"id"`
	OrderID     string         `json:"order_id"`
	Items       []PickItem     `json:"items"`
	WarehouseID string         `json:"warehouse_id"`
	Status      PickListStatus `json:"status"`
	AssignedTo  string         `json:"assigned_to,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
}

type Packing struct {
	ID          string         `json:"id"`
	PickListID  string         `json:"pick_list_id"`
	PackageID   string         `json:"package_id"`
	Weight      float64        `json:"weight"`
	Dimensions  string         `json:"dimensions"`
	Status      PackingStatus  `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
}

type Shipment struct {
	ID             string          `json:"id"`
	PackingID      string          `json:"packing_id"`
	Carrier        string          `json:"carrier"`
	TrackingNumber string          `json:"tracking_number"`
	Status         ShipmentStatus  `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
}
