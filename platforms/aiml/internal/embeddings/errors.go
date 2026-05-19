package embeddings

import "errors"

var (
	ErrEmbeddingNotFound = errors.New("embeddings: embedding not found")
	ErrInvalidDimension  = errors.New("embeddings: invalid vector dimension")
)
