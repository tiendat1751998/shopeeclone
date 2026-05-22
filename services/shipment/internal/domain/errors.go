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

	// QR code errors
	ErrQRCodeNotFound          = errors.New("QR code not found")
	ErrQRCodeExpired           = errors.New("QR code expired")
	ErrQRCodeAlreadyScanned    = errors.New("QR code already scanned")
	ErrQRCodeRevoked           = errors.New("QR code revoked")
	ErrQRCodeInactive          = errors.New("QR code inactive")
	ErrQRCodeInvalidSignature  = errors.New("QR code invalid signature")
	ErrQRCodeGenerationFailed  = errors.New("QR code generation failed")
	ErrScanUnauthorized        = errors.New("scan unauthorized: not the assigned shipper")
	ErrScanLocationInvalid     = errors.New("scan location invalid")
	ErrInvalidQRCodeType       = errors.New("invalid QR code type")
)
