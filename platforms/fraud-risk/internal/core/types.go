package core

import "time"

type EventType string

const (
	EventLogin    EventType = "login"
	EventPayment  EventType = "payment"
	EventOrder    EventType = "order"
	EventRegister EventType = "registration"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	UserID    string                 `json:"user_id"`
	IP        string                 `json:"ip"`
	DeviceID  string                 `json:"device_id"`
	Amount    float64                `json:"amount"`
	Currency  string                 `json:"currency"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type RiskLevel string

const (
	RiskSafe     RiskLevel = "safe"
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type Decision string

const (
	DecisionAllow  Decision = "allow"
	DecisionBlock  Decision = "block"
	DecisionReview Decision = "review"
)
