package warehouse

import "errors"

var (
	ErrWarehouseNotFound = errors.New("warehouse not found")
	ErrZoneNotFound      = errors.New("zone not found")
	ErrInvalidMovement   = errors.New("invalid inventory movement")
)
