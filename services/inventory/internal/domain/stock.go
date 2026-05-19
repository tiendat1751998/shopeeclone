package domain

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// StockStatus represents the current availability state of a stock item.
type StockStatus string

const (
	StockStatusInStock    StockStatus = "in_stock"
	StockStatusLowStock   StockStatus = "low_stock"
	StockStatusOutOfStock StockStatus = "out_of_stock"
	StockStatusReserved   StockStatus = "reserved"
)

func (s StockStatus) String() string { return string(s) }

func (s StockStatus) Value() (driver.Value, error) { return string(s), nil }

func (s *StockStatus) Scan(value interface{}) error {
	if value == nil { return nil }
	switch v := value.(type) {
	case string: *s = StockStatus(v)
	case []byte: *s = StockStatus(string(v))
	default: return fmt.Errorf("cannot scan %T into StockStatus", value)
	}
	return nil
}

// Stock represents inventory for a specific SKU in a specific warehouse.
// All mutations must go through Reserve/Release/Deduct/Replenish methods
// to maintain consistency between Quantity, ReservedQty, and AvailableQty.
type Stock struct {
	ID            string      `db:"id" json:"id"`
	ProductID     string      `db:"product_id" json:"product_id"`
	SkuID         string      `db:"sku_id" json:"sku_id"`
	WarehouseID   string      `db:"warehouse_id" json:"warehouse_id"`
	Quantity      int         `db:"quantity" json:"quantity"`
	ReservedQty   int         `db:"reserved_qty" json:"reserved_qty"`
	AvailableQty  int         `db:"available_qty" json:"available_qty"`
	Status        StockStatus `db:"status" json:"status"`
	ReorderLevel  int         `db:"reorder_level" json:"reorder_level"`
	Version       int         `db:"version" json:"version"` // Optimistic locking
	CreatedAt     time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at" json:"updated_at"`
}

func NewStock(productID, skuID, warehouseID string, quantity, reorderLevel int) *Stock {
	now := time.Now().UTC()
	status := StockStatusInStock
	if quantity <= 0 {
		status = StockStatusOutOfStock
	} else if quantity <= reorderLevel {
		status = StockStatusLowStock
	}
	return &Stock{
		ID: uuid.New().String(), ProductID: productID, SkuID: skuID, WarehouseID: warehouseID,
		Quantity: quantity, ReservedQty: 0, AvailableQty: quantity, Status: status,
		ReorderLevel: reorderLevel, Version: 1, CreatedAt: now, UpdatedAt: now,
	}
}

// Reserve decreases AvailableQty and increases ReservedQty.
// Must be called within a DB transaction that also creates the Reservation record.
func (s *Stock) Reserve(qty int) error {
	if qty <= 0 {
		return fmt.Errorf("%w: quantity must be positive, got %d", ErrInvalidStockOperation, qty)
	}
	if s.AvailableQty < qty {
		return ErrInsufficientStock
	}
	s.ReservedQty += qty
	s.AvailableQty -= qty
	s.UpdatedAt = time.Now().UTC()
	s.updateStatus()
	return nil
}

// Release decreases ReservedQty and increases AvailableQty.
// Must be called within a DB transaction.
func (s *Stock) Release(qty int) error {
	if qty <= 0 {
		return fmt.Errorf("%w: quantity must be positive, got %d", ErrInvalidStockOperation, qty)
	}
	if s.ReservedQty < qty {
		return fmt.Errorf("%w: cannot release %d, only %d reserved", ErrInvalidStockOperation, qty, s.ReservedQty)
	}
	s.ReservedQty -= qty
	s.AvailableQty += qty
	s.UpdatedAt = time.Now().UTC()
	s.updateStatus()
	return nil
}

// Deduct permanently removes stock (after reservation is confirmed).
// Must be called within a DB transaction.
func (s *Stock) Deduct(qty int) error {
	if qty <= 0 {
		return fmt.Errorf("%w: quantity must be positive, got %d", ErrInvalidStockOperation, qty)
	}
	if s.Quantity < qty {
		return ErrInsufficientStock
	}
	s.Quantity -= qty
	s.ReservedQty -= qty
	s.AvailableQty = s.Quantity - s.ReservedQty
	s.UpdatedAt = time.Now().UTC()
	s.updateStatus()
	return nil
}

// Replenish adds stock (e.g., return from supplier).
func (s *Stock) Replenish(qty int) {
	s.Quantity += qty
	s.AvailableQty = s.Quantity - s.ReservedQty
	s.UpdatedAt = time.Now().UTC()
	s.updateStatus()
}

func (s *Stock) updateStatus() {
	if s.AvailableQty <= 0 {
		s.Status = StockStatusOutOfStock
	} else if s.AvailableQty <= s.ReorderLevel {
		s.Status = StockStatusLowStock
	} else {
		s.Status = StockStatusInStock
	}
}

// ReservationStatus represents the lifecycle state of a stock reservation.
type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusCommitted ReservationStatus = "committed"
	ReservationStatusReleased  ReservationStatus = "released"
	ReservationStatusExpired   ReservationStatus = "expired"
)

// Reservation represents a temporary hold on stock for an order.
type Reservation struct {
	ID             string            `db:"id" json:"id"`
	OrderID        string            `db:"order_id" json:"order_id"`
	UserID         string            `db:"user_id" json:"user_id"`
	ProductID      string            `db:"product_id" json:"product_id"`
	SkuID          string            `db:"sku_id" json:"sku_id"`
	WarehouseID    string            `db:"warehouse_id" json:"warehouse_id"`
	Quantity       int               `db:"quantity" json:"quantity"`
	Status         ReservationStatus `db:"status" json:"status"`
	ExpiresAt      time.Time         `db:"expires_at" json:"expires_at"`
	IdempotencyKey string            `db:"idempotency_key" json:"idempotency_key"`
	CreatedAt      time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time         `db:"updated_at" json:"updated_at"`
}

func NewReservation(orderID, userID, productID, skuID, warehouseID string, quantity int, ttl time.Duration, idempotencyKey string) *Reservation {
	now := time.Now().UTC()
	return &Reservation{
		ID: uuid.New().String(), OrderID: orderID, UserID: userID, ProductID: productID,
		SkuID: skuID, WarehouseID: warehouseID, Quantity: quantity, Status: ReservationStatusActive,
		ExpiresAt: now.Add(ttl), IdempotencyKey: idempotencyKey, CreatedAt: now, UpdatedAt: now,
	}
}

func (r *Reservation) IsExpired() bool { return time.Now().UTC().After(r.ExpiresAt) }

func (r *Reservation) Commit() error {
	if r.Status != ReservationStatusActive {
		return fmt.Errorf("cannot commit reservation in status %s", r.Status)
	}
	r.Status = ReservationStatusCommitted
	r.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *Reservation) Release() error {
	if r.Status != ReservationStatusActive {
		return fmt.Errorf("cannot release reservation in status %s", r.Status)
	}
	r.Status = ReservationStatusReleased
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// IdempotencyRecord prevents duplicate operations.
type IdempotencyRecord struct {
	Key           string    `db:"key" json:"key"`
	ReservationID string    `db:"reservation_id" json:"reservation_id"`
	ExpiresAt     time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

func NewIdempotencyRecord(reservationID string, ttl time.Duration) *IdempotencyRecord {
	now := time.Now().UTC()
	return &IdempotencyRecord{
		Key: uuid.New().String(), ReservationID: reservationID,
		ExpiresAt: now.Add(ttl), CreatedAt: now,
	}
}

// Domain errors - exported for use in application layer error handling.
var (
	ErrStockNotFound          = fmt.Errorf("stock not found")
	ErrInsufficientStock      = fmt.Errorf("insufficient stock")
	ErrReservationNotFound    = fmt.Errorf("reservation not found")
	ErrReservationExpired     = fmt.Errorf("reservation expired")
	ErrReservationExists      = fmt.Errorf("reservation already exists")
	ErrInvalidStockOperation  = fmt.Errorf("invalid stock operation")
	ErrWarehouseNotFound      = fmt.Errorf("warehouse not found")
	ErrOversellPrevented      = fmt.Errorf("oversell prevented")
	ErrUnauthorized           = fmt.Errorf("unauthorized access")
	ErrIdempotencyKeyExists   = fmt.Errorf("idempotency key already exists")
	ErrConcurrentModification = fmt.Errorf("concurrent modification detected")
)
