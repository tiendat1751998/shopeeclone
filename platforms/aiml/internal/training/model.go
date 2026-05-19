package training

import "time"

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type TrainingJob struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	ModelName       string            `json:"model_name"`
	Dataset         string            `json:"dataset"`
	Hyperparameters map[string]string `json:"hyperparameters"`
	Status          JobStatus         `json:"status"`
	Metrics         map[string]float64 `json:"metrics,omitempty"`
	StartedAt       *time.Time        `json:"started_at,omitempty"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	Error           string            `json:"error,omitempty"`
}
