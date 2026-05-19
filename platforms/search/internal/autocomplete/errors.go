package autocomplete

import "errors"

var (
	ErrEmptyPrefix    = errors.New("prefix cannot be empty")
	ErrNoSuggestions  = errors.New("no suggestions found")
	ErrPrefixTooLong  = errors.New("prefix exceeds maximum length")
)
