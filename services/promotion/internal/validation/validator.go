package validation

import (
	"fmt"
	"strings"
)

func ValidateVoucherCode(code string) error {
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("voucher code cannot be empty")
	}
	if len(code) > 50 {
		return fmt.Errorf("voucher code too long (max 50 chars)")
	}
	return nil
}

func ValidateDiscountValue(value int64) error {
	if value <= 0 {
		return fmt.Errorf("discount value must be positive")
	}
	return nil
}

func ValidateMinSpend(spend int64) error {
	if spend < 0 {
		return fmt.Errorf("min spend cannot be negative")
	}
	return nil
}
