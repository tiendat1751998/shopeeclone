package validation
import ("fmt")
func ValidateScore(s float64) error { if s < 0 || s > 1 { return fmt.Errorf("score must be 0-1") }; return nil
func ValidateUserID(id string) error { if len(id) == 0 { return fmt.Errorf("user ID required") }; return nil }
