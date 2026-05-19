package session

import "time"

type Session struct {
	SessionID   string        `json:"session_id"`
	UserID      string        `json:"user_id"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time,omitempty"`
	Duration    time.Duration `json:"duration"`
	Pageviews   int64         `json:"pageviews"`
	EventsCount int64         `json:"events_count"`
	Device      string        `json:"device,omitempty"`
	Source      string        `json:"source,omitempty"`
	Country     string        `json:"country,omitempty"`
	IsActive    bool          `json:"is_active"`
	HasConversion bool        `json:"has_conversion"`
	Revenue     float64       `json:"revenue"`
	Events      []SessionEvent `json:"events,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}

type SessionEvent struct {
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type SessionMetrics struct {
	TotalSessions      int64   `json:"total_sessions"`
	ActiveSessions     int64   `json:"active_sessions"`
	AvgDuration        float64 `json:"avg_duration_seconds"`
	AvgPageviews       float64 `json:"avg_pageviews"`
	BounceRate         float64 `json:"bounce_rate"`
	ConversionRate     float64 `json:"conversion_rate"`
	TotalRevenue       float64 `json:"total_revenue"`
	AvgSessionRevenue  float64 `json:"avg_session_revenue"`
}

type SessionFilter struct {
	UserID    string    `json:"user_id,omitempty"`
	Source    string    `json:"source,omitempty"`
	Device    string    `json:"device,omitempty"`
	Country   string    `json:"country,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	MinDuration int64   `json:"min_duration,omitempty"`
	HasConversion *bool `json:"has_conversion,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}
