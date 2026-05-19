package query

import "errors"

var (
	ErrEmptyQuery     = errors.New("query cannot be empty")
	ErrQueryTooLong   = errors.New("query exceeds maximum length")
	ErrNoCorrections  = errors.New("no corrections available")
	ErrInvalidToken   = errors.New("invalid token")
)
