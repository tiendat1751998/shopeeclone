package inference

import "errors"

var (
	ErrModelNotFound = errors.New("inference: model not found")
	ErrInvalidInput  = errors.New("inference: invalid input")
)
