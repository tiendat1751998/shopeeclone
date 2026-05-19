package ledger

import "errors"

type AccountType string
type EntryType string
type TransactionStatus string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeEquity    AccountType = "equity"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeExpense   AccountType = "expense"

	EntryDebit  EntryType = "debit"
	EntryCredit EntryType = "credit"

	TxnStatusPending  TransactionStatus = "pending"
	TxnStatusPosted   TransactionStatus = "posted"
	TxnStatusVoided   TransactionStatus = "voided"
)

type Account struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Type     AccountType `json:"type"`
	Code     string      `json:"code"`
	Balance  int64       `json:"balance"`
	Currency string      `json:"currency"`
}

type Transaction struct {
	ID            string           `json:"id"`
	ReferenceType string           `json:"reference_type"`
	ReferenceID   string           `json:"reference_id"`
	Entries       []LedgerEntry    `json:"entries"`
	Status        TransactionStatus `json:"status"`
	CreatedAt     string           `json:"created_at"`
}

type LedgerEntry struct {
	ID            string `json:"id"`
	AccountID     string `json:"account_id"`
	TransactionID string `json:"transaction_id"`
	DebitAmount   int64  `json:"debit_amount"`
	CreditAmount  int64  `json:"credit_amount"`
	EntryType     EntryType `json:"entry_type"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
}

var (
	ErrAccountNotFound      = errors.New("ledger: account not found")
	ErrInsufficientBalance  = errors.New("ledger: insufficient balance")
	ErrInvalidAmount        = errors.New("ledger: invalid amount")
	ErrTransactionNotFound  = errors.New("ledger: transaction not found")
	ErrImbalancedEntry      = errors.New("ledger: debit and credit must balance")
	ErrAccountAlreadyExists = errors.New("ledger: account already exists")
)
