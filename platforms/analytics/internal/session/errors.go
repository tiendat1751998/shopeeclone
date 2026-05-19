package session

import "errors"

var (
	ErrSessionNotFound = errors.New("session: session not found")
	ErrSessionInvalid  = errors.New("session: invalid session data")
)
