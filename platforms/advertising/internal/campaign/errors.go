package campaign

import "errors"

var (
	ErrNotFound         = errors.New("campaign not found")
	ErrInvalidStatus    = errors.New("invalid campaign status transition")
	ErrInvalidType      = errors.New("invalid campaign type")
	ErrEmptyName        = errors.New("campaign name is required")
	ErrInvalidBudget    = errors.New("invalid campaign budget")
	ErrInvalidDateRange = errors.New("invalid campaign date range")
)
