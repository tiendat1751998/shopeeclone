package domain

import "errors"

var (
	ErrStockNotFound          = errors.New("stock not found")
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrReservationNotFound    = errors.New("reservation not found")
	ErrReservationExpired    = errors.New("reservation expired")
	ErrReservationExists     = errors.New("reservation already exists")
	ErrInvalidStockOperation = errors.New("invalid stock operation")
	ErrWarehouseNotFound      = errors.New("warehouse not found")
	ErrOversellPrevented     = errors.New("oversell prevented")
	ErrUnauthorized           = errors.New("unauthorized access")
	ErrIdempotencyKeyExists   = errors.New("idempotency key already exists")
	ErrConcurrentModification = errors.New("concurrent modification detected")
)
