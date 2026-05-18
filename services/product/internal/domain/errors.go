package domain

import (
	"fmt"
)

// DomainError is a structured domain-level error with a machine-readable code.
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError creates a new DomainError.
func NewDomainError(message, code string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// --- Not-found errors ---

// ErrProductNotFound is returned when a product cannot be located.
var ErrProductNotFound = &DomainError{
	Code:    "PRODUCT_NOT_FOUND",
	Message: "product not found",
}

// ErrSKUNotFound is returned when an SKU cannot be located.
var ErrSKUNotFound = &DomainError{
	Code:    "SKU_NOT_FOUND",
	Message: "SKU not found",
}

// ErrCategoryNotFound is returned when a category cannot be located.
var ErrCategoryNotFound = &DomainError{
	Code:    "CATEGORY_NOT_FOUND",
	Message: "category not found",
}

// ErrAttributeNotFound is returned when an attribute cannot be located.
var ErrAttributeNotFound = &DomainError{
	Code:    "ATTRIBUTE_NOT_FOUND",
	Message: "attribute not found",
}

// --- Duplicate / conflict errors ---

// ErrDuplicateProduct is returned when attempting to create a product that already exists.
var ErrDuplicateProduct = &DomainError{
	Code:    "DUPLICATE_PRODUCT",
	Message: "product already exists",
}

// ErrDuplicateSKU is returned when attempting to create an SKU that already exists.
var ErrDuplicateSKU = &DomainError{
	Code:    "DUPLICATE_SKU",
	Message: "SKU already exists",
}

// --- State errors ---

// ErrProductNotActive is returned when an operation requires an active product.
var ErrProductNotActive = &DomainError{
	Code:    "PRODUCT_NOT_ACTIVE",
	Message: "product is not active",
}

// ErrProductLocked is returned when a product is locked for editing.
var ErrProductLocked = &DomainError{
	Code:    "PRODUCT_LOCKED",
	Message: "product is locked and cannot be modified",
}

// --- Validation errors ---

// ErrInvalidPrice is returned when a price value is invalid.
var ErrInvalidPrice = &DomainError{
	Code:    "INVALID_PRICE",
	Message: "price must be greater than zero",
}

// ErrInvalidStock is returned when a stock value is invalid.
var ErrInvalidStock = &DomainError{
	Code:    "INVALID_STOCK",
	Message: "stock must be non-negative",
}

// ErrInvalidCategory is returned when a category reference is invalid.
var ErrInvalidCategory = &DomainError{
	Code:    "INVALID_CATEGORY",
	Message: "invalid category reference",
}

// --- Authorization errors ---

// ErrUnauthorizedOperation is returned when the caller lacks permission.
var ErrUnauthorizedOperation = &DomainError{
	Code:    "UNAUTHORIZED",
	Message: "unauthorized operation",
}

// IsDomainError checks if an error is a DomainError with the given code.
func IsDomainError(err error, code string) bool {
	if de, ok := err.(*DomainError); ok {
		return de.Code == code
	}
	return false
}

// IsNotFound returns true for any not-found domain error.
func IsNotFound(err error) bool {
	if de, ok := err.(*DomainError); ok {
		return de.Code == "PRODUCT_NOT_FOUND" ||
			de.Code == "SKU_NOT_FOUND" ||
			de.Code == "CATEGORY_NOT_FOUND" ||
			de.Code == "ATTRIBUTE_NOT_FOUND"
	}
	return false
}
