package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WebhookEvent struct {
	ID             string    `db:"id" json:"id"`
	PSPProvider    string    `db:"psp_provider" json:"psp_provider"`
	EventType      string    `db:"event_type" json:"event_type"`
	Payload        []byte    `db:"payload" json:"payload"`
	Signature      string    `db:"signature" json:"signature"`
	Processed      bool      `db:"processed" json:"processed"`
	RetryCount     int       `db:"retry_count" json:"retry_count"`
	IdempotencyKey string    `db:"idempotency_key" json:"idempotency_key"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

func NewWebhookEvent(pspProvider, eventType string, payload []byte, signature, idempotencyKey string) *WebhookEvent {
	return &WebhookEvent{
		ID:             uuid.New().String(),
		PSPProvider:    pspProvider,
		EventType:      eventType,
		Payload:        payload,
		Signature:      signature,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now().UTC(),
	}
}

func VerifyWebhookSignature(payload []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

type WebhookPayload struct {
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}
