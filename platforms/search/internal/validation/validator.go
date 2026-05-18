package validation
import ("fmt"; "strings")
func ValidateSearchQuery(q string) error { if len(q) > 200 { return fmt.Errorf("query too long (max 200 chars)") }; return nil }
func ValidatePrefix(p string) error { if strings.TrimSpace(p) == "" { return fmt.Errorf("prefix cannot be empty") }; if len(p) > 100 { return fmt.Errorf("prefix too long") }; return nil }
