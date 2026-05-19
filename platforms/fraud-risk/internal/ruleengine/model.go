package ruleengine

import "time"

type Rule struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	ConditionExpr   string    `json:"condition_expression"`
	Priority        int       `json:"priority"`
	Weight          float64   `json:"weight"`
	CooldownSeconds int       `json:"cooldown_seconds"`
	IsActive        bool      `json:"is_active"`
	ScoreDelta      float64   `json:"score_delta"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Strategy string

const (
	StrategyMatchAll    Strategy = "match_all"
	StrategyMatchAny    Strategy = "match_any"
	StrategyWeightedSum Strategy = "weighted_sum"
)

type RuleSet struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Rules     []Rule    `json:"rules"`
	Strategy  Strategy  `json:"strategy"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RuleEvaluation struct {
	RuleID    string  `json:"rule_id"`
	RuleName  string  `json:"rule_name"`
	Triggered bool    `json:"triggered"`
	Score     float64 `json:"score"`
	Reason    string  `json:"reason,omitempty"`
}

type RuleSetEvaluation struct {
	RuleSetID   string           `json:"ruleset_id"`
	RuleSetName string           `json:"ruleset_name"`
	Strategy    Strategy         `json:"strategy"`
	Evaluations []RuleEvaluation `json:"evaluations"`
	TotalScore  float64          `json:"total_score"`
	Passed      bool             `json:"passed"`
}
