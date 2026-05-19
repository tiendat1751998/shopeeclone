package streaming

import "time"

type StreamEvent struct {
	EntityID string    `json:"entity_id"`
	EntityType string  `json:"entity_type"`
	EventType string   `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Value     float64  `json:"value"`
}

type EventWindow struct {
	WindowSize  time.Duration `json:"window_size"`
	Count       int           `json:"count"`
	Threshold   int           `json:"threshold"`
	IsBurst     bool          `json:"is_burst"`
	AvgCount    float64       `json:"avg_count"`
}

type AggregatedCount struct {
	EntityID    string    `json:"entity_id"`
	EntityType  string    `json:"entity_type"`
	Window1Min  int       `json:"window_1min"`
	Window5Min  int       `json:"window_5min"`
	Window1Hour int       `json:"window_1hour"`
	LastUpdated time.Time `json:"last_updated"`
}
