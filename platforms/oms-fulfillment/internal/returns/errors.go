package returns

import "errors"

var (
	ErrReturnNotFound        = errors.New("return not found")
	ErrInvalidReturnData     = errors.New("invalid return data")
	ErrInvalidReturnStatus   = errors.New("invalid return status transition")
	ErrReturnAlreadyProcessed = errors.New("return already processed")
)
