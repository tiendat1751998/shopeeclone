package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReconciliationType string

const (
	ReconciliationTypePayment   ReconciliationType = "payment"
	ReconciliationTypeInventory ReconciliationType = "inventory"
	ReconciliationTypeShipment  ReconciliationType = "shipment"
)

type ReconciliationStatus string

const (
	ReconciliationStatusPending    ReconciliationStatus = "pending"
	ReconciliationStatusInProgress ReconciliationStatus = "in_progress"
	ReconciliationStatusMatched    ReconciliationStatus = "matched"
	ReconciliationStatusMismatch   ReconciliationStatus = "mismatch"
	ReconciliationStatusFailed     ReconciliationStatus = "failed"
)

type OrderReconciliation struct {
	ID                string               `db:"id" json:"id"`
	OrderID           string               `db:"order_id" json:"order_id"`
	ReconciliationType ReconciliationType  `db:"reconciliation_type" json:"reconciliation_type"`
	Status            ReconciliationStatus `db:"status" json:"status"`
	LastCheckedAt     *time.Time           `db:"last_checked_at" json:"last_checked_at,omitempty"`
	RetryCount        int                  `db:"retry_count" json:"retry_count"`
	MaxRetries        int                  `db:"max_retries" json:"max_retries"`
	Metadata          []byte               `db:"metadata" json:"metadata,omitempty"`
	CreatedAt         time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time            `db:"updated_at" json:"updated_at"`
}

func NewOrderReconciliation(orderID string, rtype ReconciliationType) *OrderReconciliation {
	now := time.Now().UTC()
	return &OrderReconciliation{
		ID:                 uuid.New().String(),
		OrderID:            orderID,
		ReconciliationType: rtype,
		Status:             ReconciliationStatusPending,
		RetryCount:         0,
		MaxRetries:         3,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (r *OrderReconciliation) CanRetry() bool {
	return r.RetryCount < r.MaxRetries
}

func (r *OrderReconciliation) IncrementRetry() {
	r.RetryCount++
	r.UpdatedAt = time.Now().UTC()
}

func (r *OrderReconciliation) MarkChecked() {
	now := time.Now().UTC()
	r.LastCheckedAt = &now
	r.UpdatedAt = now
}
