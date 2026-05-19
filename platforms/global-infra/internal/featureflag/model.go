package featureflag

import "time"

type Segment string

const (
	SegmentBetaUsers Segment = "beta_users"
	SegmentStaff     Segment = "staff"
	SegmentInternal  Segment = "internal"
	SegmentAll       Segment = "all"
)

type Rule struct {
	Attribute string `json:"attribute"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

type FeatureFlag struct {
	Name             string    `json:"name"`
	Enabled          bool      `json:"enabled"`
	PercentageRollout int      `json:"percentage_rollout"`
	UserSegment      Segment   `json:"user_segment"`
	Rules            []Rule    `json:"rules,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type EvaluateRequest struct {
	FlagName string `json:"flag_name" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
}

type EvaluateResponse struct {
	FlagName string `json:"flag_name"`
	UserID   string `json:"user_id"`
	Enabled  bool   `json:"enabled"`
	Reason   string `json:"reason"`
}
