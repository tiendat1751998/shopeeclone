package validation
import ("fmt"; "strings")
func ValidateProductName(n string) error { if strings.TrimSpace(n) == "" { return fmt.Errorf("product name cannot be empty") }; if len(n) > 500 { return fmt.Errorf("product name too long") }; return nil }
func ValidateCategoryID(id string) error { if strings.TrimSpace(id) == "" { return fmt.Errorf("category ID cannot be empty") }; return nil }
func ValidatePrice(p int64) error { if p < 0 { return fmt.Errorf("price cannot be negative") }; return nil }
