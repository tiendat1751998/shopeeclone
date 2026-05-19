package dispatch

import "errors"

var (
	ErrDispatchNotFound     = errors.New("dispatch not found")
	ErrCourierNotAvailable  = errors.New("courier not available")
	ErrDispatchAlreadyAssigned = errors.New("dispatch already assigned")
	ErrInvalidDispatchStatus  = errors.New("invalid dispatch status")
)
