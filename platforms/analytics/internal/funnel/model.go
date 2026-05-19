package funnel

import "time"

type FunnelStep struct {
	Name     string    `json:"name"`
	EventType string   `json:"event_type"`
	Order    int       `json:"order"`
	Window   string    `json:"window,omitempty"`
}

type FunnelDefinition struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Steps       []FunnelStep `json:"steps"`
	TimeRange   string       `json:"time_range"`
	CreatedAt   time.Time    `json:"created_at"`
}

type FunnelResult struct {
	ID             string            `json:"id"`
	FunnelName     string            `json:"funnel_name"`
	Steps          []FunnelStepResult `json:"steps"`
	OverallRate    float64           `json:"overall_conversion_rate"`
	StartCount     int64             `json:"start_count"`
	EndCount       int64             `json:"end_count"`
	AnalyzedAt     time.Time         `json:"analyzed_at"`
}

type FunnelStepResult struct {
	StepName      string  `json:"step_name"`
	EventType     string  `json:"event_type"`
	Order         int     `json:"order"`
	UserCount     int64   `json:"user_count"`
	StepRate      float64 `json:"step_conversion_rate"`
	OverallRate   float64 `json:"overall_conversion_rate"`
	DropOff       int64   `json:"drop_off_count"`
	DropOffRate   float64 `json:"drop_off_rate"`
}

type ConversionRate struct {
	StepFrom string  `json:"step_from"`
	StepTo   string  `json:"step_to"`
	Rate     float64 `json:"rate"`
	FromCount int64  `json:"from_count"`
	ToCount   int64  `json:"to_count"`
}
