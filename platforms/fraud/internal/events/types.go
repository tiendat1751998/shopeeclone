package events

import "time"

type EventType string

const (
	EventFraudAlertTriggered EventType = "fraud.alert.triggered"
	EventRuleUpdated         EventType = "fraud.rule.updated"
	EventCaseCreated         EventType = "fraud.case.created"
	EventBlacklistHit        EventType = "fraud.blacklist.hit"
)

type FraudAlertTriggeredEvent struct {
	AlertID   string  `json:"alert_id"`
	EventID   string  `json:"event_id"`
	UserID    string  `json:"user_id"`
	AlertType string  `json:"alert_type"`
	RiskScore float64 `json:"risk_score"`
	RiskLevel string  `json:"risk_level"`
	Timestamp time.Time `json:"timestamp"`
}

type RuleUpdatedEvent struct {
	RuleID    string `json:"rule_id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CaseCreatedEvent struct {
	CaseID    string  `json:"case_id"`
	AlertID   string  `json:"alert_id"`
	UserID    string  `json:"user_id"`
	Priority  string  `json:"priority"`
	RiskScore float64 `json:"risk_score"`
	CreatedAt time.Time `json:"created_at"`
}

type BlacklistHitEvent struct {
	EntityType string `json:"entity_type"`
	EntityValue string `json:"entity_value"`
	Reason     string `json:"reason"`
	Timestamp  time.Time `json:"timestamp"`
}
