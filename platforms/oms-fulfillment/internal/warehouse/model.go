package warehouse

import "time"

type ZoneType string

const (
	ZoneTypeStorage ZoneType = "storage"
	ZoneTypePicking ZoneType = "picking"
	ZoneTypePacking ZoneType = "packing"
	ZoneTypeShipping ZoneType = "shipping"
)

type MovementType string

const (
	MovementTransfer MovementType = "transfer"
	MovementReceive  MovementType = "receive"
	MovementShip     MovementType = "ship"
	MovementAdjust   MovementType = "adjust"
)

type Warehouse struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Address           string `json:"address"`
	City              string `json:"city"`
	State             string `json:"state"`
	IsActive          bool   `json:"is_active"`
	Capacity          int    `json:"capacity"`
	CurrentUtilization int   `json:"current_utilization"`
}

type Zone struct {
	ID          string   `json:"id"`
	WarehouseID string   `json:"warehouse_id"`
	Name        string   `json:"name"`
	Type        ZoneType `json:"type"`
}

type InventoryMovement struct {
	ID          string       `json:"id"`
	ProductID   string       `json:"product_id"`
	WarehouseID string       `json:"warehouse_id"`
	FromZone    string       `json:"from_zone,omitempty"`
	ToZone      string       `json:"to_zone,omitempty"`
	Quantity    int          `json:"quantity"`
	Type        MovementType `json:"type"`
	Reference   string       `json:"reference,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}
