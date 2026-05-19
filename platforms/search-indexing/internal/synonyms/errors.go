package synonyms

import "errors"

var (
	ErrSynonymSetNotFound = errors.New("synonym set not found")
	ErrSetAlreadyExists   = errors.New("synonym set already exists")
	ErrEmptyWords         = errors.New("synonym set must contain at least one word")
)
