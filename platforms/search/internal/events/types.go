package events

import "time"

type EventType string

const (
	EventDocumentIndexed  EventType = "document.indexed"
	EventDocumentDeleted  EventType = "document.deleted"
	EventSearchPerformed  EventType = "search.performed"
	EventIndexTaskCreated EventType = "index.task.created"
)

type DocumentIndexedEvent struct {
	DocumentID  string    `json:"document_id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	SellerID    string    `json:"seller_id"`
	IndexedAt   time.Time `json:"indexed_at"`
}

type DocumentDeletedEvent struct {
	DocumentID string    `json:"document_id"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type SearchPerformedEvent struct {
	Query       string `json:"query"`
	ResultCount int64  `json:"result_count"`
	TookMs      int64  `json:"took_ms"`
	UserID      string `json:"user_id,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
}

type IndexTaskCreatedEvent struct {
	TaskID         string    `json:"task_id"`
	DocumentID     string    `json:"document_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
}
