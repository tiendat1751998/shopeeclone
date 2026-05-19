package reconciliation

import "errors"

type RunStatus string
type ItemStatus string

const (
	RunStatusPending    RunStatus = "pending"
	RunStatusInProgress RunStatus = "in_progress"
	RunStatusCompleted  RunStatus = "completed"
	RunStatusFailed     RunStatus = "failed"

	ItemStatusMatched     ItemStatus = "matched"
	ItemStatusUnmatched   ItemStatus = "unmatched"
	ItemStatusDiscrepancy ItemStatus = "discrepancy"
)

type ReconciliationRun struct {
	ID                string    `json:"id"`
	Date              string    `json:"date"`
	Status            RunStatus `json:"status"`
	TotalTransactions int       `json:"total_transactions"`
	MatchedCount      int       `json:"matched_count"`
	UnmatchedCount    int       `json:"unmatched_count"`
	DiscrepancyAmount int64     `json:"discrepancy_amount"`
}

type ReconciliationItem struct {
	RunID         string     `json:"run_id"`
	TransactionID string     `json:"transaction_id"`
	SourceAmount  int64      `json:"source_amount"`
	LedgerAmount  int64      `json:"ledger_amount"`
	Difference    int64      `json:"difference"`
	Status        ItemStatus `json:"status"`
	Notes         string     `json:"notes,omitempty"`
}

var (
	ErrRunNotFound = errors.New("reconciliation: run not found")
)
