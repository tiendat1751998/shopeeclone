package validation
import ("fmt")
func ValidateBudget(b int64) error { if b <= 0 { return fmt.Errorf("budget must be positive") }; return nil
func ValidateBidAmount(b int64) error { if b < 0 { return fmt.Errorf("bid cannot be negative") }; return nil
