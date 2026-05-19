package vectorstore

import "errors"

var (
	ErrRecordNotFound   = errors.New("vectorstore: record not found")
	ErrInvalidDimension = errors.New("vectorstore: invalid vector dimension")
	ErrEmptyVector      = errors.New("vectorstore: empty vector not allowed")
)
