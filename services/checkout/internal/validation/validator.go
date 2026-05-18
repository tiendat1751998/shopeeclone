package validation

import (
	"fmt"
	"strings"
)

func ValidateCheckoutID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("checkout ID cannot be empty")
	}
	return nil
}

func ValidateUserID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	return nil
}

func ValidateCartID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("cart ID cannot be empty")
	}
	return nil
}

func ValidateGrandTotal(total int64) error {
	if total < 0 {
		return fmt.Errorf("grand total cannot be negative")
	}
	return nil
}
