package fraudcase

import "time"

type CaseStatus string

const (
	StatusOpen       CaseStatus = "open"
	StatusInvestigating CaseStatus = "investigating"
	StatusEscalated  CaseStatus = "escalated"
	StatusResolved   CaseStatus = "resolved"
	StatusClosed     CaseStatus = "closed"
)

type CasePriority string

const (
	PriorityLow      CasePriority = "low"
	PriorityMedium   CasePriority = "medium"
	PriorityHigh     CasePriority = "high"
	PriorityCritical CasePriority = "critical"
)

type FraudCase struct {
	ID            string          `json:"id"`
	AlertID       string          `json:"alert_id"`
	UserID        string          `json:"user_id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Status        CaseStatus      `json:"status"`
	Priority      CasePriority    `json:"priority"`
	RiskScore     float64         `json:"risk_score"`
	Investigator  string          `json:"investigator,omitempty"`
	Evidence      []Evidence      `json:"evidence,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	ResolvedAt    *time.Time      `json:"resolved_at,omitempty"`
	Resolution    string          `json:"resolution,omitempty"`
}

type Evidence struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Data        string    `json:"data"`
	AddedBy     string    `json:"added_by"`
	AddedAt     time.Time `json:"added_at"`
}

type Attachment struct {
	ID       string `json:"id"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	Data     []byte `json:"-"`
	URL      string `json:"url,omitempty"`
	Size     int64  `json:"size"`
}
