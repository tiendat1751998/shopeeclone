package detection

import "errors"

var (
	ErrAlertNotFound = errors.New("fraud: alert not found")
	ErrEventInvalid  = errors.New("fraud: invalid event")
	ErrAlertResolved = errors.New("fraud: alert already resolved")
)
