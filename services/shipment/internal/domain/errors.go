package domain

import "errors"

var (
	ErrShipmentNotFound         = errors.New("shipment not found")
	ErrShipmentAlreadyExists    = errors.New("shipment already exists")
	ErrInvalidShipmentState     = errors.New("invalid shipment state transition")
	ErrCarrierUnavailable       = errors.New("carrier unavailable")
	ErrInvalidWebhookSignature  = errors.New("invalid webhook signature")
	ErrWebhookReplayDetected    = errors.New("webhook replay detected")
	ErrTrackingNotFound         = errors.New("tracking not found")
	ErrUnauthorized             = errors.New("unauthorized access")
	ErrIdempotencyKeyExists     = errors.New("idempotency key already exists")
	ErrConcurrentModification   = errors.New("concurrent modification detected")
)
