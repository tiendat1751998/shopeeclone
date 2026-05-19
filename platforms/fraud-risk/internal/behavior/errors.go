package behavior

import "errors"

var (
	ErrProfileNotFound = errors.New("fraud-risk: behavior profile not found")
	ErrRuleNotFound    = errors.New("fraud-risk: behavioral rule not found")
)
