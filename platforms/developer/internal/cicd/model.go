package cicd

import "time"

type TriggerType string

const (
	TriggerPush     TriggerType = "push"
	TriggerPR       TriggerType = "pr"
	TriggerSchedule TriggerType = "schedule"
)

type PipelineStatus string

const (
	StatusPending PipelineStatus = "pending"
	StatusRunning PipelineStatus = "running"
	StatusSuccess PipelineStatus = "success"
	StatusFailed  PipelineStatus = "failed"
)

type StageStatus string

const (
	StagePending StageStatus = "pending"
	StageRunning StageStatus = "running"
	StageSuccess StageStatus = "success"
	StageFailed  StageStatus = "failed"
)

type Stage struct {
	Name            string       `json:"name"`
	Status          StageStatus  `json:"status"`
	DurationSeconds int          `json:"duration_seconds"`
	Logs            string       `json:"logs"`
}

type Pipeline struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Service     string         `json:"service"`
	Trigger     TriggerType    `json:"trigger"`
	Stages      []Stage        `json:"stages"`
	Status      PipelineStatus `json:"status"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CommitSHA   string         `json:"commit_sha"`
}
