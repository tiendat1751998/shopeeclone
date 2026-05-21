package validation
import ("fmt"; "strings")
func ValidateEventType(t string) error {
	valid := map[string]bool{"page_view":true,"product_view":true,"click":true,"add_to_cart":true,"checkout":true,"search":true,"impression":true}
	if !valid[t] { return fmt.Errorf("invalid event type: %s", t) }; return nil
}
func ValidateUserID(id string) error { if strings.TrimSpace(id) == "" { return fmt.Errorf("user ID cannot be empty") }; return nil }
