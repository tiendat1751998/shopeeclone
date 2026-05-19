package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Reservation represents a temporary hold on stock for an order.
// NOTE: This type is also defined in stock.go - this file is kept for backward compatibility.
// Use the Reservation type from stock.go for new code.

type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusCommitted ReservationStatus = "committed"
	ReservationStatusReleased  ReservationStatus = "released"
	ReservationStatusExpired   ReservationStatus = "expired"
)

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
