package experiments

import "errors"

var (
	ErrExperimentNotFound = errors.New("experiments: experiment not found")
	ErrExperimentExists   = errors.New("experiments: experiment already exists")
	ErrExperimentClosed   = errors.New("experiments: experiment is not running")
)
