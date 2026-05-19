package inventory

import "time"

type ReservationStatus string

const (
	ReservationPending   ReservationStatus = "pending"
	ReservationReserved  ReservationStatus = "reserved"
	ReservationReleased  ReservationStatus = "released"
	ReservationConsumed  ReservationStatus = "consumed"
)

type InventoryReservation struct {
	ID        string            `json:"id"`
	OrderID   string            `json:"order_id"`
	ProductID string            `json:"product_id"`
	SKU       string            `json:"sku"`
	Quantity  int               `json:"quantity"`
	Status    ReservationStatus `json:"status"`
	ExpiresAt time.Time         `json:"expires_at"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Stock struct {
	ProductID   string `json:"product_id"`
	WarehouseID string `json:"warehouse_id"`
	Available   int    `json:"available"`
	Reserved    int    `json:"reserved"`
	Total       int    `json:"total"`
}

type ReserveRequest struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
}
