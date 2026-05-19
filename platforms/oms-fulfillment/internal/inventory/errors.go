package inventory

import "errors"

var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrInsufficientStock   = errors.New("insufficient stock available")
	ErrInvalidReservation  = errors.New("invalid reservation data")
	ErrStockNotFound       = errors.New("stock not found")
)
