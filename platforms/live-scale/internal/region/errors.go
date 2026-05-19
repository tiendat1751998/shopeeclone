package region

import "errors"

var (
	ErrRegionNotFound      = errors.New("region not found")
	ErrNoRegionAvailable   = errors.New("no region available")
	ErrInvalidLatencyData  = errors.New("invalid latency data")
)
