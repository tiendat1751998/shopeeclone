package validation
import ("fmt"; "strings")
func ValidateTitle(t string) error { if strings.TrimSpace(t) == "" { return fmt.Errorf("title cannot be empty") }; return nil }
func ValidateContent(c string) error { if len(c) > 500 { return fmt.Errorf("content too long") }; return nil }
