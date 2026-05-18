package domain

import (
	"time"

	"github.com/google/uuid"
)

type IdempotencyRecord struct {
	Key         string    `db:"key" json:"key"`
	ShipmentID  string    `db:"shipment_id" json:"shipment_id"`
	ExpiresAt   time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func NewIdempotencyRecord(shipmentID string, ttl time.Duration) *IdempotencyRecord {
	now := time.Now().UTC()
	return &IdempotencyRecord{
		Key: uuid.New().String(), ShipmentID: shipmentID, ExpiresAt: now.Add(ttl), CreatedAt: now,
	}
}
