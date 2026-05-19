package search

import "errors"

var (
	ErrNotFound      = errors.New("document not found")
	ErrInvalidQuery  = errors.New("invalid search query")
	ErrIndexFailed   = errors.New("index operation failed")
	ErrSearchFailed  = errors.New("search operation failed")
	ErrNoResults     = errors.New("no results found")
	ErrInvalidFacet  = errors.New("invalid facet field")
	ErrInvalidSort   = errors.New("invalid sort field")
)
