package realtime

import "time"

type SessionEvent struct {
	EventType string    `json:"event_type"`
	ItemID    string    `json:"item_id,omitempty"`
	Query     string    `json:"query,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type UserSession struct {
	UserID       string         `json:"user_id"`
	SessionID    string         `json:"session_id"`
	Events       []SessionEvent `json:"events"`
	StartedAt    time.Time      `json:"started_at"`
	LastActiveAt time.Time      `json:"last_active_at"`
}

type ArmStat struct {
	ArmID    string  `json:"arm_id"`
	Plays    int     `json:"plays"`
	Rewards  float64 `json:"rewards"`
	MeanReward float64 `json:"mean_reward"`
}
