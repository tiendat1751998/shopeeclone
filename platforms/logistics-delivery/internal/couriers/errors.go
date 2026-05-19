package couriers

import "errors"

var (
	ErrCourierNotFound        = errors.New("courier not found")
	ErrCourierNotAvailable    = errors.New("courier not available")
	ErrInvalidWebhookSignature = errors.New("invalid webhook signature")
	ErrDuplicateWebhookEvent   = errors.New("duplicate webhook event")
	ErrCourierAtMaxCapacity    = errors.New("courier at max capacity")
)
