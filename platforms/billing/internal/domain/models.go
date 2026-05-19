package domain

import (
	"time"
	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeEquity    AccountType = "equity"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeExpense   AccountType = "expense"
)

type Account struct {
	ID        string      `db:"id" json:"id"`
	UserID    string      `db:"user_id" json:"user_id"`
	Type      AccountType `db:"type" json:"type"`
	Currency  string      `db:"currency" json:"currency"`
	Balance   int64       `db:"balance" json:"balance"`
	Frozen    int64       `db:"frozen" json:"frozen"`
	Status    string      `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}

func NewAccount(userID string, accType AccountType, currency string) *Account {
	return &Account{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      accType,
		Currency:  currency,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type LedgerEntry struct {
	ID            string    `db:"id" json:"id"`
	TransactionID string    `db:"transaction_id" json:"transaction_id"`
	AccountID     string    `db:"account_id" json:"account_id"`
	Type          EntryType `db:"type" json:"type"`
	Amount        int64     `db:"amount" json:"amount"`
	Currency      string    `db:"currency" json:"currency"`
	BalanceBefore int64     `db:"balance_before" json:"balance_before"`
	BalanceAfter  int64     `db:"balance_after" json:"balance_after"`
	Description   string    `db:"description" json:"description"`
	Reference     string    `db:"reference" json:"reference"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type EntryType string

const (
	EntryDebit  EntryType = "debit"
	EntryCredit EntryType = "credit"
)

type Transaction struct {
	ID            string          `db:"id" json:"id"`
	IDempotencyKey string         `db:"idempotency_key" json:"idempotency_key"`
	Type          TransactionType `db:"type" json:"type"`
	Status        TransactionStatus `db:"status" json:"status"`
	Amount        int64           `db:"amount" json:"amount"`
	Currency      string          `db:"currency" json:"currency"`
	Description   string          `db:"description" json:"description"`
	DebitAccountID  string        `db:"debit_account_id" json:"debit_account_id"`
	CreditAccountID string        `db:"credit_account_id" json:"credit_account_id"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	CompletedAt   *time.Time      `db:"completed_at" json:"completed_at"`
}

type TransactionType string

const (
	TxnDeposit    TransactionType = "deposit"
	TxnWithdrawal TransactionType = "withdrawal"
	TxnTransfer   TransactionType = "transfer"
	TxnPayment    TransactionType = "payment"
	TxnRefund     TransactionType = "refund"
	TxnFee        TransactionType = "fee"
	TxnSettlement TransactionType = "settlement"
	TxnPayout     TransactionType = "payout"
	TxnAdjustment TransactionType = "adjustment"
)

type TransactionStatus string

const (
	TxnPending    TransactionStatus = "pending"
	TxnCompleted  TransactionStatus = "completed"
	TxnFailed     TransactionStatus = "failed"
	TxnReversed   TransactionStatus = "reversed"
	TxnExpired    TransactionStatus = "expired"
)

type Wallet struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Type      string    `db:"type" json:"type"`
	Currency  string    `db:"currency" json:"currency"`
	Balance   int64     `db:"balance" json:"balance"`
	Frozen    int64     `db:"frozen" json:"frozen"`
	Pending   int64     `db:"pending" json:"pending"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Settlement struct {
	ID            string          `db:"id" json:"id"`
	MerchantID    string          `db:"merchant_id" json:"merchant_id"`
	Amount        int64           `db:"amount" json:"amount"`
	Currency      string          `db:"currency" json:"currency"`
	FeeAmount     int64           `db:"fee_amount" json:"fee_amount"`
	NetAmount     int64           `db:"net_amount" json:"net_amount"`
	Status        SettlementStatus `db:"status" json:"status"`
	PeriodStart   time.Time       `db:"period_start" json:"period_start"`
	PeriodEnd     time.Time       `db:"period_end" json:"period_end"`
	ScheduledDate *time.Time      `db:"scheduled_date" json:"scheduled_date"`
	CompletedAt   *time.Time      `db:"completed_at" json:"completed_at"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}

type SettlementStatus string

const (
	SettlementPending   SettlementStatus = "pending"
	SettlementProcessing SettlementStatus = "processing"
	SettlementCompleted SettlementStatus = "completed"
	SettlementFailed    SettlementStatus = "failed"
)

type Payout struct {
	ID          string        `db:"id" json:"id"`
	MerchantID  string        `db:"merchant_id" json:"merchant_id"`
	Amount      int64         `db:"amount" json:"amount"`
	Currency    string        `db:"currency" json:"currency"`
	Method      string        `db:"method" json:"method"`
	Status      PayoutStatus  `db:"status" json:"status"`
	AccountRef  string        `db:"account_ref" json:"account_ref"`
	Description string        `db:"description" json:"description"`
	RequestedAt time.Time     `db:"requested_at" json:"requested_at"`
	CompletedAt *time.Time    `db:"completed_at" json:"completed_at"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
}

type PayoutStatus string

const (
	PayoutRequested  PayoutStatus = "requested"
	PayoutProcessing PayoutStatus = "processing"
	PayoutCompleted  PayoutStatus = "completed"
	PayoutFailed     PayoutStatus = "failed"
)

type Refund struct {
	ID              string        `db:"id" json:"id"`
	TransactionID   string        `db:"transaction_id" json:"transaction_id"`
	OriginalTxnID   string        `db:"original_txn_id" json:"original_txn_id"`
	Amount          int64         `db:"amount" json:"amount"`
	Currency        string        `db:"currency" json:"currency"`
	Reason          string        `db:"reason" json:"reason"`
	Status          RefundStatus  `db:"status" json:"status"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
}

type RefundStatus string

const (
	RefundPending    RefundStatus = "pending"
	RefundProcessing RefundStatus = "processing"
	RefundCompleted  RefundStatus = "completed"
	RefundFailed     RefundStatus = "failed"
)

type ReconciliationRun struct {
	ID           string        `db:"id" json:"id"`
	PeriodStart  time.Time     `db:"period_start" json:"period_start"`
	PeriodEnd    time.Time     `db:"period_end" json:"period_end"`
	Status       string        `db:"status" json:"status"`
	TotalMatched int64         `db:"total_matched" json:"total_matched"`
	TotalMismatch int64        `db:"total_mismatch" json:"total_mismatch"`
	CreatedAt    time.Time     `db:"created_at" json:"created_at"`
}

type AuditLog struct {
	ID          string    `db:"id" json:"id"`
	ActorID     string    `db:"actor_id" json:"actor_id"`
	Action      string    `db:"action" json:"action"`
	Resource    string    `db:"resource" json:"resource"`
	ResourceID  string    `db:"resource_id" json:"resource_id"`
	Details     string    `db:"details" json:"details"`
	IPAddress   string    `db:"ip_address" json:"ip_address"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

var (
	ErrInsufficientBalance = ErrBilling("insufficient_balance")
	ErrAccountNotFound     = ErrBilling("account_not_found")
	ErrWalletNotFound      = ErrBilling("wallet_not_found")
	ErrDuplicateTxn        = ErrBilling("duplicate_transaction")
	ErrInvalidAmount       = ErrBilling("invalid_amount")
	ErrAccountFrozen       = ErrBilling("account_frozen")
	ErrCurrencyMismatch    = ErrBilling("currency_mismatch")
	ErrLedgerImbalance     = ErrBilling("ledger_imbalance")
	ErrRefundAlreadyExists = ErrBilling("refund_already_exists")
	ErrSettlementFailed    = ErrBilling("settlement_failed")
)

type ErrBilling string
func (e ErrBilling) Error() string { return "billing: " + string(e) }
