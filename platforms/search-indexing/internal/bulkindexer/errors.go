package bulkindexer

import "errors"

var (
	ErrJobNotFound      = errors.New("bulk job not found")
	ErrBatchNotFound    = errors.New("document batch not found")
	ErrJobAlreadyExists = errors.New("bulk job already exists")
	ErrJobCompleted     = errors.New("job already completed")
)
