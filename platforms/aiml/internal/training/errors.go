package training

import "errors"

var (
	ErrJobNotFound = errors.New("training: job not found")
	ErrJobExists   = errors.New("training: job already exists")
	ErrInvalidTransition = errors.New("training: invalid status transition")
)
