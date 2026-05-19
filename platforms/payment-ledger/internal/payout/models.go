package payout

import "errors"

type PayoutStatus string
type PaymentMethod string

const (
	StatusPending    PayoutStatus = "pending"
	StatusProcessing PayoutStatus = "processing"
	StatusCompleted  PayoutStatus = "completed"
	StatusFailed     PayoutStatus = "failed"

	MethodBankTransfer PaymentMethod = "bank_transfer"
	MethodWallet       PaymentMethod = "wallet"
)

type Payout struct {
	ID            string        `json:"id"`
	SellerID      string        `json:"seller_id"`
	Amount        int64         `json:"amount"`
	Fee           int64         `json:"fee"`
	NetAmount     int64         `json:"net_amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Status        PayoutStatus  `json:"status"`
	BatchID       string        `json:"batch_id,omitempty"`
	PeriodStart   string        `json:"period_start"`
	PeriodEnd     string        `json:"period_end"`
	CompletedAt   string        `json:"completed_at,omitempty"`
	CreatedAt     string        `json:"created_at"`
}

type PayoutBatch struct {
	ID          string        `json:"id"`
	TotalAmount int64         `json:"total_amount"`
	Count       int           `json:"count"`
	Status      PayoutStatus  `json:"status"`
	CreatedAt   string        `json:"created_at"`
}

var (
	ErrPayoutNotFound     = errors.New("payout: not found")
	ErrInvalidPayoutAmount = errors.New("payout: invalid amount")
	ErrInvalidPayoutStatus = errors.New("payout: invalid status transition")
)
