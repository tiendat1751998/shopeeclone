package validation
import ("fmt"; "strings")
func ValidateChannel(c string) error {
	valid := map[string]bool{"push": true, "email": true, "sms": true, "inapp": true}
	if !valid[c] { return fmt.Errorf("invalid channel: %s", c) }; return nil
}
func ValidateTitle(t string) error { if strings.TrimSpace(t) == "" { return fmt.Errorf("title cannot be empty") }; if len(t) > 500 { return fmt.Errorf("title too long") }; return nil
}
