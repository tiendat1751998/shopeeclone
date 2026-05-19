package routing

import "errors"

var (
	ErrRouteNotFound       = errors.New("route not found")
	ErrZoneNotFound        = errors.New("zone not found")
	ErrNoAvailableCourier  = errors.New("no available courier for zone")
	ErrNoWarehouseInZone   = errors.New("no warehouse in zone")
	ErrInvalidWaypointData = errors.New("invalid waypoint data")
)
