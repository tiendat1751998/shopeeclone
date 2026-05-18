package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReconciliationStatus string

const (
	ReconStatusPending    ReconciliationStatus = "pending"
	ReconStatusMatched    ReconciliationStatus = "matched"
	ReconStatusMismatch   ReconciliationStatus = "mismatch"
	ReconStatusFailed     ReconciliationStatus = "failed"
)

type PaymentReconciliation struct {
	ID            string               `db:"id" json:"id"`
	PaymentID     string               `db:"payment_id" json:"payment_id"`
	OrderID       string               `db:"order_id" json:"order_id"`
	Type          string               `db:"type" json:"type"`
	Status        ReconciliationStatus `db:"status" json:"status"`
	PSPReference  string               `db:"psp_reference" json:"psp_reference"`
	Amount        int64                `db:"amount" json:"amount"`
	Discrepancy   int64                `db:"discrepancy" json:"discrepancy"`
	RetryCount    int                  `db:"retry_count" json:"retry_count"`
	LastCheckedAt *time.Time           `db:"last_checked_at" json:"last_checked_at,omitempty"`
	Metadata      []byte               `db:"metadata" json:"metadata,omitempty"`
	CreatedAt     time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `db:"updated_at" json:"updated_at"`
}

func NewPaymentReconciliation(paymentID, orderID, rtype string, amount int64) *PaymentReconciliation {
	now := time.Now().UTC()
	return &PaymentReconciliation{
		ID:        uuid.New().String(),
		PaymentID: paymentID,
		OrderID:   orderID,
		Type:      rtype,
		Status:    ReconStatusPending,
		Amount:    amount,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
