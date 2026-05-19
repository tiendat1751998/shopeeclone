package similarity

import "errors"

var (
	ErrEmptyQueryVector = errors.New("similarity: empty query vector not allowed")
)
