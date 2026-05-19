package behavior

import "time"

type UserBehaviorProfile struct {
	UserID           string    `json:"user_id"`
	TypicalLoginHour int       `json:"typical_login_hour"`
	TypicalIPRange   string    `json:"typical_ip_range"`
	TypicalDevice    string    `json:"typical_device"`
	ActionSequence   []string  `json:"action_sequence"`
	LastUpdated      time.Time `json:"last_updated"`
}

type BehavioralRule struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Condition string `json:"condition"`
	Threshold int    `json:"threshold"`
}

type DeviationResult struct {
	HasDeviation bool     `json:"has_deviation"`
	Reasons      []string `json:"reasons"`
	Score        float64  `json:"score"`
}
