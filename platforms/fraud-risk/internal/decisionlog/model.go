package decisionlog

import "time"

type DecisionType string

const (
	DecisionAllow  DecisionType = "allow"
	DecisionBlock  DecisionType = "block"
	DecisionReview DecisionType = "review"
)

type DecisionLog struct {
	ID            string       `json:"id"`
	EventID       string       `json:"event_id"`
	EventType     string       `json:"event_type"`
	UserID        string       `json:"user_id"`
	Decision      DecisionType `json:"decision"`
	RiskScore     float64      `json:"risk_score"`
	TriggeredRules []string    `json:"triggered_rules"`
	Timestamp     time.Time    `json:"timestamp"`
}

type DecisionStats struct {
	TotalDecisions int            `json:"total_decisions"`
	AllowCount     int            `json:"allow_count"`
	BlockCount     int            `json:"block_count"`
	ReviewCount    int            `json:"review_count"`
	AvgRiskScore   float64        `json:"avg_risk_score"`
}
