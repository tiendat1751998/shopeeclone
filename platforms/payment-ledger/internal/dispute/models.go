package dispute

import "errors"

type DisputeStatus string
type Resolution string

const (
	StatusOpened     DisputeStatus = "opened"
	StatusUnderReview DisputeStatus = "under_review"
	StatusResolved   DisputeStatus = "resolved"
	StatusClosed     DisputeStatus = "closed"

	ResolutionFullRefund    Resolution = "full_refund"
	ResolutionPartialRefund Resolution = "partial_refund"
	ResolutionNoRefund      Resolution = "no_refund"
)

type Dispute struct {
	ID            string        `json:"id"`
	TransactionID string        `json:"transaction_id"`
	PaymentID     string        `json:"payment_id"`
	UserID        string        `json:"user_id"`
	Reason        string        `json:"reason"`
	Amount        int64         `json:"amount"`
	Status        DisputeStatus `json:"status"`
	Resolution    Resolution    `json:"resolution,omitempty"`
	Evidence      []string      `json:"evidence,omitempty"`
	OpenedAt      string        `json:"opened_at"`
	ResolvedAt    string        `json:"resolved_at,omitempty"`
}

var (
	ErrDisputeNotFound   = errors.New("dispute: not found")
	ErrInvalidDisputeStatus = errors.New("dispute: invalid status transition")
	ErrAlreadyResolved   = errors.New("dispute: already resolved")
)
