package payment

import "errors"

type PaymentStatus string
type PaymentMethod string

const (
	StatusPending            PaymentStatus = "pending"
	StatusProcessing         PaymentStatus = "processing"
	StatusCompleted          PaymentStatus = "completed"
	StatusFailed             PaymentStatus = "failed"
	StatusRefunded           PaymentStatus = "refunded"
	StatusPartiallyRefunded  PaymentStatus = "partially_refunded"

	MethodCreditCard PaymentMethod = "credit_card"
	MethodBankTransfer PaymentMethod = "bank_transfer"
	MethodWallet      PaymentMethod = "wallet"
	MethodCOD         PaymentMethod = "cod"
)

type Payment struct {
	ID              string        `json:"id"`
	OrderID         string        `json:"order_id"`
	UserID          string        `json:"user_id"`
	Amount          int64         `json:"amount"`
	Currency        string        `json:"currency"`
	Method          PaymentMethod `json:"method"`
	Status          PaymentStatus `json:"status"`
	GatewayResponse string        `json:"gateway_response,omitempty"`
	TransactionID   string        `json:"transaction_id,omitempty"`
	Fee             int64         `json:"fee"`
	NetAmount       int64         `json:"net_amount"`
	AuthorizedAt    string        `json:"authorized_at,omitempty"`
	CapturedAt      string        `json:"captured_at,omitempty"`
	SettledAt       string        `json:"settled_at,omitempty"`
	CreatedAt       string        `json:"created_at"`
	UpdatedAt       string        `json:"updated_at"`
}

var (
	ErrPaymentNotFound    = errors.New("payment: not found")
	ErrInvalidAmount      = errors.New("payment: invalid amount")
	ErrInvalidMethod      = errors.New("payment: invalid method")
	ErrInvalidStatus      = errors.New("payment: invalid status transition")
	ErrAlreadyRefunded    = errors.New("payment: already refunded")
	ErrRefundExceeds      = errors.New("payment: refund exceeds paid amount")
)
