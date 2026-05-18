package domain

import "errors"

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInvalidProductData   = errors.New("invalid product data")
	ErrSKUNotFound          = errors.New("SKU not found")
	ErrSKUAlreadyExists     = errors.New("SKU already exists")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryHasChildren  = errors.New("category has child categories")
	ErrCategoryHasProducts  = errors.New("category has products")
	ErrAttributeNotFound    = errors.New("attribute not found")
	ErrMediaNotFound        = errors.New("media not found")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrInsufficientPerms    = errors.New("insufficient permissions")
	ErrIdempotencyKeyExists = errors.New("idempotency key already exists")
	ErrInvalidParentCategory = errors.New("invalid parent category")
	ErrCircularCategory     = errors.New("circular category reference")
	ErrConcurrentModification = errors.New("concurrent modification detected")
	ErrIndexingFailed       = errors.New("search indexing failed")
	ErrModerationRequired   = errors.New("product requires moderation")
)
