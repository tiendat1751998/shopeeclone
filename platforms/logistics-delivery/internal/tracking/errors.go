package tracking

import "errors"

var (
	ErrTrackingEventNotFound = errors.New("tracking event not found")
	ErrDuplicateEvent        = errors.New("duplicate tracking event")
	ErrInvalidEventType      = errors.New("invalid tracking event type")
	ErrShipmentNotTrackable  = errors.New("shipment not trackable")
)
