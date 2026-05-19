package ordermanagement

import "errors"

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrInvalidOrderData       = errors.New("invalid order data")
	ErrEmptyItems             = errors.New("order must have at least one item")
)
