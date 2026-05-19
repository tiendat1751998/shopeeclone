package cdn

import "errors"

var (
	ErrEndpointNotFound   = errors.New("cdn endpoint not found")
	ErrNoEndpointsAvailable = errors.New("no cdn endpoints available")
	ErrPurgeFailed        = errors.New("cdn purge failed")
	ErrInvalidPurgeRequest = errors.New("invalid purge request")
)
