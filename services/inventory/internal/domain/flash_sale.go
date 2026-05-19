package domain

import "time"

type FlashSaleInventory struct {
	ID            string    `db:"id" json:"id"`
	FlashSaleID   string    `db:"flash_sale_id" json:"flash_sale_id"`
	SKU           string    `db:"sku" json:"sku"`
	WarehouseID   string    `db:"warehouse_id" json:"warehouse_id"`
	TotalStock    int       `db:"total_stock" json:"total_stock"`
	ReservedStock int       `db:"reserved_stock" json:"reserved_stock"`
	SoldStock     int       `db:"sold_stock" json:"sold_stock"`
	MaxPerUser    int       `db:"max_per_user" json:"max_per_user"`
	StartTime     time.Time `db:"start_time" json:"start_time"`
	EndTime       time.Time `db:"end_time" json:"end_time"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (fs *FlashSaleInventory) IsAvailable() bool {
	if !fs.IsActive {
		return false
	}
	now := time.Now()
	return now.After(fs.StartTime) && now.Before(fs.EndTime)
}

func (fs *FlashSaleInventory) AvailableStock() int {
	return fs.TotalStock - fs.ReservedStock - fs.SoldStock
}
