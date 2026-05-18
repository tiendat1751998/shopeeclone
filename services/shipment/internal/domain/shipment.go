package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus string

const (
	ShipmentStatusPending    ShipmentStatus = "pending"
	ShipmentStatusBooked     ShipmentStatus = "booked"
	ShipmentStatusPickedUp   ShipmentStatus = "picked_up"
	ShipmentStatusInTransit  ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered  ShipmentStatus = "delivered"
	ShipmentStatusFailed     ShipmentStatus = "failed"
	ShipmentStatusReturned   ShipmentStatus = "returned"
	ShipmentStatusCancelled  ShipmentStatus = "cancelled"
)

func (s ShipmentStatus) String() string { return string(s) }
func (s ShipmentStatus) Value() (driver.Value, error) { return string(s), nil }
func (s *ShipmentStatus) Scan(value interface{}) error {
	if value == nil { return nil }
	switch v := value.(type) {
	case string: *s = ShipmentStatus(v)
	case []byte: *s = ShipmentStatus(string(v))
	default: return fmt.Errorf("cannot scan %T", value)
	}
	return nil
}

type Shipment struct {
	ID               string         `db:"id" json:"id"`
	OrderID          string         `db:"order_id" json:"order_id"`
	UserID           string         `db:"user_id" json:"user_id"`
	CarrierID        string         `db:"carrier_id" json:"carrier_id"`
	TrackingNumber   string         `db:"tracking_number" json:"tracking_number"`
	Status           ShipmentStatus `db:"status" json:"status"`
	OriginAddress    Address        `db:"origin_address" json:"origin_address"`
	DestAddress      Address        `db:"destination_address" json:"destination_address"`
	Weight           float64        `db:"weight" json:"weight"`
	Dimensions       string         `db:"dimensions" json:"dimensions"`
	LabelURL         string         `db:"label_url" json:"label_url"`
	Cost             int64          `db:"cost" json:"cost"`
	Currency         string         `db:"currency" json:"currency"`
	IdempotencyKey   string         `db:"idempotency_key" json:"idempotency_key"`
	Metadata         json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	Version          int            `db:"version" json:"version"`
	EstimatedDelivery *time.Time    `db:"estimated_delivery" json:"estimated_delivery,omitempty"`
	DeliveredAt      *time.Time     `db:"delivered_at" json:"delivered_at,omitempty"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updated_at"`
}

type Address struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Phone      string `json:"phone,omitempty"`
}

func NewShipment(orderID, userID, carrierID, idempotencyKey, currency string, origin, dest Address, weight float64) *Shipment {
	now := time.Now().UTC()
	return &Shipment{
		ID: uuid.New().String(), OrderID: orderID, UserID: userID, CarrierID: carrierID,
		Status: ShipmentStatusPending, OriginAddress: origin, DestAddress: dest,
		Weight: weight, Currency: currency, IdempotencyKey: idempotencyKey,
		Version: 1, CreatedAt: now, UpdatedAt: now,
	}
}

func (s *Shipment) CanTransitionTo(target ShipmentStatus) bool {
	valid := map[ShipmentStatus][]ShipmentStatus{
		ShipmentStatusPending:    {ShipmentStatusBooked, ShipmentStatusCancelled},
		ShipmentStatusBooked:     {ShipmentStatusPickedUp, ShipmentStatusCancelled},
		ShipmentStatusPickedUp:   {ShipmentStatusInTransit, ShipmentStatusFailed},
		ShipmentStatusInTransit:  {ShipmentStatusOutForDelivery, ShipmentStatusFailed},
		ShipmentStatusOutForDelivery: {ShipmentStatusDelivered, ShipmentStatusFailed},
		ShipmentStatusDelivered:  {ShipmentStatusReturned},
		ShipmentStatusFailed:     {ShipmentStatusInTransit, ShipmentStatusReturned},
		ShipmentStatusReturned:   {},
		ShipmentStatusCancelled:  {},
	}
	allowed, ok := valid[s.Status]; if !ok { return false }
	for _, st := range allowed { if st == target { return true } }
	return false
}

func (s *Shipment) TransitionTo(target ShipmentStatus) error {
	if !s.CanTransitionTo(target) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidShipmentState, s.Status, target)
	}
	now := time.Now().UTC()
	s.Status = target; s.Version++; s.UpdatedAt = now
	if target == ShipmentStatusDelivered { s.DeliveredAt = &now }
	return nil
}

func (s *Shipment) IsTerminal() bool {
	return s.Status == ShipmentStatusDelivered || s.Status == ShipmentStatusReturned || s.Status == ShipmentStatusCancelled
}
