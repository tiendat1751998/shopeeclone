package pipeline

import "errors"

var (
	ErrPipelineNotFound    = errors.New("pipeline not found")
	ErrPipelineExists     = errors.New("pipeline already exists")
	ErrInvalidStageConfig = errors.New("invalid stage config")
	ErrUnknownProcessor   = errors.New("unknown processor type")
)
