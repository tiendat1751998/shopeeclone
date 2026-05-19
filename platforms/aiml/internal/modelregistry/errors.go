package modelregistry

import "errors"

var (
	ErrModelNotFound    = errors.New("modelregistry: model not found")
	ErrModelExists     = errors.New("modelregistry: model already exists")
	ErrInvalidStage    = errors.New("modelregistry: invalid stage transition")
)
