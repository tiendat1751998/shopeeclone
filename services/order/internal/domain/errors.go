package domain

import "errors"

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrOrderAlreadyExists      = errors.New("order already exists")
	ErrInvalidStateTransition  = errors.New("invalid order state transition")
	ErrOrderNotCancellable     = errors.New("order cannot be cancelled in current state")
	ErrOrderNotModifiable      = errors.New("order cannot be modified")
	ErrUnauthorized            = errors.New("unauthorized access to order")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrIdempotencyKeyExists    = errors.New("idempotency key already exists")
	ErrInvalidOrderData        = errors.New("invalid order data")
	ErrSellerSplitFailed       = errors.New("seller order split failed")
	ErrReconciliationFailed    = errors.New("order reconciliation failed")
	ErrSnapshotMismatch        = errors.New("order snapshot checksum mismatch")
	ErrDuplicateEvent          = errors.New("duplicate event detected")
	ErrInvalidIdempotencyKey   = errors.New("invalid idempotency key")
	ErrOrderExpired            = errors.New("order has expired")
	ErrPaymentRequired         = errors.New("payment required for this operation")
	ErrConcurrentModification  = errors.New("concurrent modification detected")
)
