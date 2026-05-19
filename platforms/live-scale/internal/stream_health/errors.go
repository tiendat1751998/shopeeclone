package stream_health

import "errors"

var (
	ErrStreamNotFound = errors.New("stream health not found")
	ErrInvalidMetric  = errors.New("invalid health metric")
)
