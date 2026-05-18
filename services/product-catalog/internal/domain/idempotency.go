package domain

import (
	"time"

	"github.com/google/uuid"
)

type IdempotencyRecord struct {
	Key         string    `db:"key" json:"key"`
	ProductID   string    `db:"product_id" json:"product_id"`
	ExpiresAt   time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func NewIdempotencyRecord(productID string, ttl time.Duration) *IdempotencyRecord {
	now := time.Now().UTC()
	return &IdempotencyRecord{
		Key: uuid.New().String(), ProductID: productID,
		ExpiresAt: now.Add(ttl), CreatedAt: now,
	}
}

func (r *IdempotencyRecord) IsExpired() bool {
	return time.Now().UTC().After(r.ExpiresAt)
}
