package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type IdempotencyRecord struct {
	Key       string    `db:"key" json:"key"`
	PaymentID string    `db:"payment_id" json:"payment_id"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewIdempotencyRecord(key, paymentID string, ttl time.Duration) *IdempotencyRecord {
	now := time.Now().UTC()
	return &IdempotencyRecord{
		Key:       key,
		PaymentID: paymentID,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}
}

func GenerateIdempotencyKey(orderID, userID string, amount int64, timestamp time.Time) string {
	data, _ := json.Marshal(map[string]interface{}{
		"order_id":  orderID,
		"user_id":   userID,
		"amount":    amount,
		"timestamp": timestamp.Unix(),
	})
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
