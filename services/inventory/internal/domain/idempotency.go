package domain

import (
	"time"

	"github.com/google/uuid"
)

type IdempotencyRecord struct {
	Key           string    `db:"key" json:"key"`
	ReservationID string    `db:"reservation_id" json:"reservation_id"`
	ExpiresAt     time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

func NewIdempotencyRecord(reservationID string, ttl time.Duration) *IdempotencyRecord {
	now := time.Now().UTC()
	return &IdempotencyRecord{Key: uuid.New().String(), ReservationID: reservationID, ExpiresAt: now.Add(ttl), CreatedAt: now}
}
