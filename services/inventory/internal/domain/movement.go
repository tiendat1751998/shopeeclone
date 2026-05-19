package domain

import "time"

type StockMovement struct {
	ID          string    `db:"id" json:"id"`
	SKU         string    `db:"sku" json:"sku"`
	WarehouseID string    `db:"warehouse_id" json:"warehouse_id"`
	Type        string    `db:"type" json:"type"`
	Quantity    int       `db:"quantity" json:"quantity"`
	BeforeQty   int       `db:"before_qty" json:"before_qty"`
	AfterQty    int       `db:"after_qty" json:"after_qty"`
	ReferenceID string    `db:"reference_id" json:"reference_id"`
	Reason      string    `db:"reason" json:"reason"`
	OperatorID  string    `db:"operator_id" json:"operator_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
