package bulkindexer

import "time"

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type BatchStatus string

const (
	BatchStatusPending   BatchStatus = "pending"
	BatchStatusProcessed BatchStatus = "processed"
	BatchStatusFailed    BatchStatus = "failed"
)

type BulkJob struct {
	ID               string      `json:"id"`
	IndexName        string      `json:"index_name"`
	TotalDocuments   int         `json:"total_documents"`
	ProcessedCount   int         `json:"processed_count"`
	FailedCount      int         `json:"failed_count"`
	Status           JobStatus   `json:"status"`
	Errors           []string    `json:"errors,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	CompletedAt      *time.Time  `json:"completed_at,omitempty"`
}

type DocumentBatch struct {
	ID           string      `json:"id"`
	JobID        string      `json:"job_id"`
	Documents    []map[string]interface{} `json:"documents"`
	BatchNumber  int         `json:"batch_number"`
	Status       BatchStatus `json:"status"`
}
