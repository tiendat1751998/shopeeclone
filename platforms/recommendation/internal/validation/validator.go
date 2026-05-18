package validation
import ("fmt"; "strings")
func ValidateUserID(id string) error { if strings.TrimSpace(id) == "" { return fmt.Errorf("user ID cannot be empty") }; return nil }
func ValidateProductID(id string) error { if strings.TrimSpace(id) == "" { return fmt.Errorf("product ID cannot be empty") }; return nil }
func ValidateLimit(limit int) error { if limit <= 0 || limit > 100 { return fmt.Errorf("limit must be 1-100") }; return nil }
