package ruleengine

import "errors"

var (
	ErrRuleNotFound    = errors.New("fraud-risk: rule not found")
	ErrRuleSetNotFound = errors.New("fraud-risk: ruleset not found")
)
