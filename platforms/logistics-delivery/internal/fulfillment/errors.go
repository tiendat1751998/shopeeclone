package fulfillment

import "errors"

var (
	ErrFulfillmentNotFound = errors.New("fulfillment not found")
	ErrAlreadyPacked       = errors.New("items already packed")
	ErrInvalidFulfillmentStatus = errors.New("invalid fulfillment status")
)
