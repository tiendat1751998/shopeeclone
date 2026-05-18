package validation

import (
	"fmt"
	"strings"
)

func ValidateSKU(sku string) error {
	if strings.TrimSpace(sku) == "" {
		return fmt.Errorf("SKU cannot be empty")
	}
	if len(sku) > 100 {
		return fmt.Errorf("SKU too long (max 100 chars)")
	}
	return nil
}

func ValidateQuantity(qty, maxQty int) error {
	if qty <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if qty > maxQty {
		return fmt.Errorf("quantity exceeds maximum of %d", maxQty)
	}
	return nil
}

func ValidateCartID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("cart ID cannot be empty")
	}
	return nil
}

func ValidateUserID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	return nil
}

func ValidateUnitPrice(price int64) error {
	if price < 0 {
		return fmt.Errorf("unit price cannot be negative")
	}
	return nil
}
