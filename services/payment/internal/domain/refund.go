package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusProcessed RefundStatus = "processed"
	RefundStatusFailed    RefundStatus = "failed"
)

type Refund struct {
	ID               string       `db:"id" json:"id"`
	PaymentID        string       `db:"payment_id" json:"payment_id"`
	OrderID          string       `db:"order_id" json:"order_id"`
	Amount           int64        `db:"amount" json:"amount"`
	Currency         string       `db:"currency" json:"currency"`
	Status           RefundStatus `db:"status" json:"status"`
	Reason           string       `db:"reason" json:"reason"`
	PSPRefundID      string       `db:"psp_refund_id" json:"psp_refund_id,omitempty"`
	IdempotencyKey   string       `db:"idempotency_key" json:"idempotency_key"`
	Metadata         []byte       `db:"metadata" json:"metadata,omitempty"`
	CreatedAt        time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time    `db:"updated_at" json:"updated_at"`
}

func NewRefund(paymentID, orderID, currency, reason, idempotencyKey string, amount int64) *Refund {
	now := time.Now().UTC()
	return &Refund{
		ID:             uuid.New().String(),
		PaymentID:      paymentID,
		OrderID:        orderID,
		Amount:         amount,
		Currency:       currency,
		Status:         RefundStatusPending,
		Reason:         reason,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
