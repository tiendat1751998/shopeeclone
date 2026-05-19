package fraudcase

import "errors"

var (
	ErrCaseNotFound = errors.New("fraud: case not found")
	ErrInvalidTransition = errors.New("fraud: invalid case status transition")
)
