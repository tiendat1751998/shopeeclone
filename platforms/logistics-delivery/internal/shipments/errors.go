package shipments

import "errors"

var (
	ErrShipmentNotFound       = errors.New("shipment not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrDuplicateReplayID       = errors.New("duplicate replay id")
	ErrShipmentAlreadyExists   = errors.New("shipment already exists")
	ErrInvalidShipmentData     = errors.New("invalid shipment data")
	ErrShipmentAlreadyDelivered = errors.New("shipment already delivered")
	ErrShipmentCancelled       = errors.New("shipment already cancelled")
)
