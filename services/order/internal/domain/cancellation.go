package domain

import (
	"time"

	"github.com/google/uuid"
)

type CancellationType string

const (
	CancellationTypeUser    CancellationType = "user"
	CancellationTypeSeller  CancellationType = "seller"
	CancellationTypeTimeout CancellationType = "timeout"
	CancellationTypeSystem  CancellationType = "system"
)

type CompensationStatus string

const (
	CompensationPending    CompensationStatus = "pending"
	CompensationInProgress CompensationStatus = "in_progress"
	CompensationCompleted  CompensationStatus = "completed"
	CompensationFailed     CompensationStatus = "failed"
	CompensationSkipped    CompensationStatus = "skipped"
)

type OrderCancellation struct {
	ID                string             `db:"id" json:"id"`
	OrderID           string             `db:"order_id" json:"order_id"`
	Reason            string             `db:"reason" json:"reason"`
	CancelledBy       string             `db:"cancelled_by" json:"cancelled_by"`
	CancelledByType   CancellationType   `db:"cancelled_by_type" json:"cancelled_by_type"`
	CompensationStatus CompensationStatus `db:"compensation_status" json:"compensation_status"`
	RefundAmount      int64              `db:"refund_amount" json:"refund_amount"`
	Metadata          []byte             `db:"metadata" json:"metadata,omitempty"`
	CreatedAt         time.Time          `db:"created_at" json:"created_at"`
}

func NewOrderCancellation(orderID, reason, cancelledBy string, cancelledByType CancellationType, refundAmount int64) *OrderCancellation {
	return &OrderCancellation{
		ID:                 uuid.New().String(),
		OrderID:            orderID,
		Reason:             reason,
		CancelledBy:        cancelledBy,
		CancelledByType:    cancelledByType,
		CompensationStatus: CompensationPending,
		RefundAmount:       refundAmount,
		CreatedAt:          time.Now().UTC(),
	}
}
