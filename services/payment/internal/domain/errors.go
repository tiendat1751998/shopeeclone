package domain

import "errors"

var (
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrPaymentAlreadyExists     = errors.New("payment already exists")
	ErrPaymentAlreadyProcessed  = errors.New("payment already processed")
	ErrInvalidPaymentState      = errors.New("invalid payment state transition")
	ErrDoubleChargeDetected     = errors.New("double charge detected")
	ErrUnauthorized             = errors.New("unauthorized access")
	ErrInsufficientPermissions  = errors.New("insufficient permissions")
	ErrInvalidWebhookSignature  = errors.New("invalid webhook signature")
	ErrWebhookReplayDetected    = errors.New("webhook replay detected")
	ErrPSPUnavailable           = errors.New("payment service provider unavailable")
	ErrRefundNotAllowed         = errors.New("refund not allowed for this payment")
	ErrRefundAmountExceeded     = errors.New("refund amount exceeds payment amount")
	ErrIdempotencyKeyExists     = errors.New("idempotency key already exists")
	ErrFraudDetected            = errors.New("fraud detected")
	ErrInvalidPaymentMethod     = errors.New("invalid payment method")
	ErrSettlementFailed         = errors.New("settlement failed")
	ErrConcurrentModification   = errors.New("concurrent modification detected")
)
