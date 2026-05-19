package transcoding

import "time"

type JobStatus string

const (
	JobPending    JobStatus = "pending"
	JobProcessing JobStatus = "processing"
	JobCompleted  JobStatus = "completed"
	JobFailed     JobStatus = "failed"
)

type VideoProfile string

const (
	Profile480p  VideoProfile = "480p"
	Profile720p  VideoProfile = "720p"
	Profile1080p VideoProfile = "1080p"
)

type VideoProfileConfig struct {
	Profile    VideoProfile `json:"profile"`
	Width      int          `json:"width"`
	Height     int          `json:"height"`
	BitrateKbps int         `json:"bitrate_kbps"`
}

var SupportedProfiles = map[VideoProfile]VideoProfileConfig{
	Profile480p:  {Profile: Profile480p, Width: 854, Height: 480, BitrateKbps: 1500},
	Profile720p:  {Profile: Profile720p, Width: 1280, Height: 720, BitrateKbps: 3000},
	Profile1080p: {Profile: Profile1080p, Width: 1920, Height: 1080, BitrateKbps: 6000},
}

type Output struct {
	Profile   VideoProfile `json:"profile"`
	URL       string       `json:"url"`
	SizeBytes int64        `json:"size_bytes"`
	DurationMs int64       `json:"duration_ms"`
}

type TranscodeJob struct {
	ID         string       `json:"id"`
	StreamID   string       `json:"stream_id"`
	InputURL   string       `json:"input_url"`
	Profiles   []VideoProfile `json:"profiles"`
	Status     JobStatus    `json:"status"`
	Outputs    []Output     `json:"outputs,omitempty"`
	Error      string       `json:"error,omitempty"`
	Progress   float64      `json:"progress"`
	CreatedAt  time.Time    `json:"created_at"`
	StartedAt  *time.Time   `json:"started_at,omitempty"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
	UpdatedAt  time.Time    `json:"updated_at"`
}
