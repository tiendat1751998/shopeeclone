package verification

import "errors"

var (
	ErrVerificationNotFound = errors.New("fraud: verification not found")
	ErrCodeMismatch         = errors.New("fraud: verification code mismatch")
	ErrCodeExpired          = errors.New("fraud: verification code expired")
	ErrMaxAttempts          = errors.New("fraud: max verification attempts reached")
)
