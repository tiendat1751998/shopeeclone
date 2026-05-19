package rules

type RuleDefinition struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Condition   string  `json:"condition"`
	Severity    int     `json:"severity"`
	Weight      float64 `json:"weight"`
	IsActive    bool    `json:"is_active"`
	Cooldown    int     `json:"cooldown"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type RuleEvaluation struct {
	RuleID    string  `json:"rule_id"`
	RuleName  string  `json:"rule_name"`
	Severity  int     `json:"severity"`
	Weight    float64 `json:"weight"`
	Score     float64 `json:"score"`
	Triggered bool    `json:"triggered"`
	Reason    string  `json:"reason,omitempty"`
}
