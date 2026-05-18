package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SKUStatus string

const (
	SKUStatusActive   SKUStatus = "active"
	SKUStatusInactive SKUStatus = "inactive"
	SKUStatusSoldOut  SKUStatus = "sold_out"
)

func (s SKUStatus) String() string { return string(s) }
func (s SKUStatus) Value() (driver.Value, error) { return string(s), nil }
func (s *SKUStatus) Scan(value interface{}) error {
	if value == nil { return nil }
	switch v := value.(type) {
	case string: *s = SKUStatus(v)
	case []byte: *s = SKUStatus(string(v))
	default: return fmt.Errorf("cannot scan %T into SKUStatus", value)
	}
	return nil
}

type SKU struct {
	ID         string          `db:"id" json:"id"`
	ProductID  string          `db:"product_id" json:"product_id"`
	SkuCode    string          `db:"sku_code" json:"sku_code"`
	Name       string          `db:"name" json:"name"`
	Price      int64           `db:"price" json:"price"`
	ComparePrice int64         `db:"compare_price" json:"compare_price"`
	Currency   string          `db:"currency" json:"currency"`
	Stock      int             `db:"stock" json:"stock"`
	ReservedStock int          `db:"reserved_stock" json:"reserved_stock"`
	Weight     float64         `db:"weight" json:"weight"`
	Dimensions string          `db:"dimensions" json:"dimensions"`
	Status     SKUStatus       `db:"status" json:"status"`
	Attributes json.RawMessage `db:"attributes" json:"attributes,omitempty"`
	Metadata   json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	Version    int             `db:"version" json:"version"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
}

func NewSKU(productID, skuCode, name, currency string, price int64) *SKU {
	now := time.Now().UTC()
	return &SKU{
		ID: uuid.New().String(), ProductID: productID, SkuCode: skuCode,
		Name: name, Price: price, Currency: currency, Status: SKUStatusActive,
		Version: 1, CreatedAt: now, UpdatedAt: now,
	}
}

func (s *SKU) IsAvailable() bool {
	return s.Status == SKUStatusActive && s.Stock > s.ReservedStock
}

func (s *SKU) AvailableStock() int {
	avail := s.Stock - s.ReservedStock
	if avail < 0 { return 0 }
	return avail
}

func (s *SKU) Reserve(quantity int) error {
	if s.AvailableStock() < quantity {
		return fmt.Errorf("%w: available %d, requested %d", ErrSKUNotFound, s.AvailableStock(), quantity)
	}
	s.ReservedStock += quantity
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *SKU) Release(quantity int) error {
	if s.ReservedStock < quantity {
		return fmt.Errorf("cannot release more than reserved")
	}
	s.ReservedStock -= quantity
	s.UpdatedAt = time.Now().UTC()
	return nil
}
