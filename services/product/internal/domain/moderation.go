package domain

import (
	"time"
)

// ModerationStatus represents the review state of a product.
type ModerationStatus string

const (
	ModerationStatusPending  ModerationStatus = "PENDING"
	ModerationStatusApproved ModerationStatus = "APPROVED"
	ModerationStatusRejected ModerationStatus = "REJECTED"
	ModerationStatusFlagged  ModerationStatus = "FLAGGED"
)

// ModerationRecord tracks a single moderation review action on a product.
type ModerationRecord struct {
	ID         string           `db:"id"          json:"id"`
	SPUID      string           `db:"spu_id"      json:"spu_id"`
	Status     ModerationStatus `db:"status"      json:"status"`
	Reason     string           `db:"reason"      json:"reason,omitempty"`
	ReviewerID string           `db:"reviewer_id" json:"reviewer_id,omitempty"`
	CreatedAt  time.Time        `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time        `db:"updated_at"  json:"updated_at"`
}

// IsResolved returns true if the moderation decision has been made.
func (mr *ModerationRecord) IsResolved() bool {
	return mr.Status == ModerationStatusApproved || mr.Status == ModerationStatusRejected
}

// IsPending returns true if the record is awaiting review.
func (mr *ModerationRecord) IsPending() bool {
	return mr.Status == ModerationStatusPending
}

// IsApproved returns true if the product was approved.
func (mr *ModerationRecord) IsApproved() bool {
	return mr.Status == ModerationStatusApproved
}

// IsRejected returns true if the product was rejected.
func (mr *ModerationRecord) IsRejected() bool {
	return mr.Status == ModerationStatusRejected
}

// IsFlagged returns true if the product was flagged for further review.
func (mr *ModerationRecord) IsFlagged() bool {
	return mr.Status == ModerationStatusFlagged
}
