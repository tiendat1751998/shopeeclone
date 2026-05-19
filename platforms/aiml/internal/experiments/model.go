package experiments

import "time"

type ExperimentStatus string

const (
	StatusRunning  ExperimentStatus = "running"
	StatusPaused   ExperimentStatus = "paused"
	StatusCompleted ExperimentStatus = "completed"
)

type MetricType string

const (
	MetricCTR      MetricType = "ctr"
	MetricRevenue  MetricType = "revenue"
	MetricConversion MetricType = "conversion"
)

type Variant struct {
	ModelName             string  `json:"model_name"`
	TrafficAllocationPct  float64 `json:"traffic_allocation_percentage"`
}

type Experiment struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	ModelA           string           `json:"model_a"`
	ModelB           string           `json:"model_b"`
	TrafficPct       float64          `json:"traffic_percentage"`
	Metric           MetricType       `json:"metric"`
	Status           ExperimentStatus `json:"status"`
	Variants         []Variant        `json:"variants"`
	StartedAt        time.Time        `json:"started_at"`
}

type Assignment struct {
	ExperimentID string `json:"experiment_id"`
	UserID       string `json:"user_id"`
	Variant      string `json:"variant"`
}

type Result struct {
	ExperimentID string    `json:"experiment_id"`
	Variant      string    `json:"variant"`
	UserID       string    `json:"user_id"`
	Value        float64   `json:"value"`
	Timestamp    time.Time `json:"timestamp"`
}

type ExperimentResults struct {
	ExperimentID string  `json:"experiment_id"`
	Name         string  `json:"name"`
	Metric       MetricType `json:"metric"`
	ModelA       string  `json:"model_a"`
	ModelB       string  `json:"model_b"`
	AResults     VariantResults `json:"a_results"`
	BResults     VariantResults `json:"b_results"`
	Improvement  float64 `json:"improvement_percentage"`
}

type VariantResults struct {
	ModelName   string   `json:"model_name"`
	SampleSize  int      `json:"sample_size"`
	Sum         float64  `json:"sum"`
	Mean        float64  `json:"mean"`
}
