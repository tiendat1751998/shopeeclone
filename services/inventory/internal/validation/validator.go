package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var skuPattern = regexp.MustCompile(`^[A-Za-z0-9\-_]{1,100}$`)

// ValidateSKU validates SKU format
func ValidateSKU(sku string) error {
	if strings.TrimSpace(sku) == "" {
		return fmt.Errorf("SKU cannot be empty")
	}
	if !skuPattern.MatchString(sku) {
		return fmt.Errorf("SKU format invalid: must be alphanumeric, hyphens, underscores (max 100 chars)")
	}
	return nil
}

// ValidateQuantity validates quantity is positive
func ValidateQuantity(qty int64) error {
	if qty <= 0 {
		return fmt.Errorf("quantity must be positive, got %d", qty)
	}
	return nil
}

// ValidateWarehouseID validates warehouse ID format
func ValidateWarehouseID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("warehouse ID cannot be empty")
	}
	if len(id) > 36 {
		return fmt.Errorf("warehouse ID too long (max 36 chars)")
	}
	return nil
}

// ValidateReservationKey validates reservation key format
func ValidateReservationKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("reservation key cannot be empty")
	}
	if len(key) > 100 {
		return fmt.Errorf("reservation key too long (max 100 chars)")
	}
	return nil
}

// ValidateUserID validates user ID format
func ValidateUserID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if len(id) > 36 {
		return fmt.Errorf("user ID too long (max 36 chars)")
	}
	return nil
}

// ValidateIdempotencyKey validates idempotency key
func ValidateIdempotencyKey(key string) error {
	if key == "" {
		return nil // optional
	}
	if len(key) > 100 {
		return fmt.Errorf("idempotency key too long (max 100 chars)")
	}
	return nil
}
