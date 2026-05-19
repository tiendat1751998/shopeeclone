package indexing

import "errors"

var (
	ErrTaskNotFound      = errors.New("index task not found")
	ErrDuplicateIdempotency = errors.New("duplicate idempotency key")
	ErrInvalidDocument   = errors.New("invalid document")
	ErrIndexingFailed    = errors.New("indexing failed")
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
)
