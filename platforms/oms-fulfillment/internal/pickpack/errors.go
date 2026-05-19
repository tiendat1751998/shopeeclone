package pickpack

import "errors"

var (
	ErrPickListNotFound   = errors.New("pick list not found")
	ErrPackingNotFound    = errors.New("packing not found")
	ErrShipmentNotFound   = errors.New("shipment not found")
	ErrInvalidPickData    = errors.New("invalid pick data")
	ErrInvalidPackData    = errors.New("invalid packing data")
	ErrInvalidShipData    = errors.New("invalid shipment data")
)
