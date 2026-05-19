package ranking

import "errors"

var (
	ErrConfigNotFound = errors.New("ranking config not found")
	ErrInvalidFactor  = errors.New("invalid ranking factor")
)
