package blacklist

import "time"

type BlacklistType string

const (
	BlacklistIP     BlacklistType = "ip"
	BlacklistUser   BlacklistType = "user"
	BlacklistCard   BlacklistType = "card"
	BlacklistDevice BlacklistType = "device"
)

type BlacklistReason string

const (
	ReasonFraudulentActivity BlacklistReason = "fraudulent_activity"
	ReasonChargeback         BlacklistReason = "chargeback"
	ReasonSuspiciousBehavior BlacklistReason = "suspicious_behavior"
	ReasonManualReview       BlacklistReason = "manual_review"
	ReasonHighRisk           BlacklistReason = "high_risk"
)

type BlacklistEntry struct {
	ID        string          `json:"id"`
	Type      BlacklistType   `json:"type"`
	Value     string          `json:"value"`
	Reason    BlacklistReason `json:"reason"`
	AddedBy   string          `json:"added_by"`
	CreatedAt time.Time       `json:"created_at"`
	ExpiresAt *time.Time      `json:"expires_at,omitempty"`
	IsActive  bool            `json:"is_active"`
}

type CheckRequest struct {
	UserID   string `json:"user_id"`
	IP       string `json:"ip"`
	DeviceID string `json:"device_id"`
	CardNumber string `json:"card_number"`
}

type CheckResponse struct {
	Blocked bool              `json:"blocked"`
	Reasons []BlacklistReason `json:"reasons,omitempty"`
	Entries []BlacklistEntry  `json:"entries,omitempty"`
}
