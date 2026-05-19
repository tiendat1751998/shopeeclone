package reconciliation

import (
	"context"

	"github.com/google/uuid"
)

type PaymentSource interface {
	GetAmount(paymentID string) int64
}

type LedgerSource interface {
	GetAmount(transactionID string) int64
}

type Service struct {
	repo          Repository
	paymentSource PaymentSource
	ledgerSource  LedgerSource
}

func NewService(repo Repository, ps PaymentSource, ls LedgerSource) *Service {
	return &Service{
		repo:          repo,
		paymentSource: ps,
		ledgerSource:  ls,
	}
}

type PaymentRecord struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
}

type LedgerRecord struct {
	TransactionID string `json:"transaction_id"`
	Amount        int64  `json:"amount"`
}

func (s *Service) RunReconciliation(ctx context.Context, date string, payments []PaymentRecord, ledgers []LedgerRecord) (*ReconciliationRun, error) {
	run := &ReconciliationRun{
		ID:     uuid.New().String(),
		Date:   date,
		Status: RunStatusPending,
	}

	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, err
	}

	run.Status = RunStatusInProgress

	pmap := make(map[string]int64)
	for _, p := range payments {
		pmap[p.ID] = p.Amount
	}
	lmap := make(map[string]int64)
	for _, l := range ledgers {
		lmap[l.TransactionID] = l.Amount
	}

	seen := make(map[string]bool)
	var items []*ReconciliationItem
	matched := 0
	unmatched := 0
	var discrepancy int64

	for _, p := range payments {
		ledgerAmt, ok := lmap[p.ID]
		seen[p.ID] = true
		diff := p.Amount - ledgerAmt
		item := &ReconciliationItem{
			RunID:         run.ID,
			TransactionID: p.ID,
			SourceAmount:  p.Amount,
			LedgerAmount:  ledgerAmt,
			Difference:    diff,
		}
		if !ok {
			item.Status = ItemStatusUnmatched
			item.Notes = "payment found in source but missing in ledger"
			unmatched++
		} else if diff != 0 {
			item.Status = ItemStatusDiscrepancy
			item.Notes = "amount mismatch between source and ledger"
			discrepancy += diff
			unmatched++
		} else {
			item.Status = ItemStatusMatched
			matched++
		}
		items = append(items, item)
	}

	for _, l := range ledgers {
		if seen[l.TransactionID] {
			continue
		}
		item := &ReconciliationItem{
			RunID:         run.ID,
			TransactionID: l.TransactionID,
			SourceAmount:  0,
			LedgerAmount:  l.Amount,
			Difference:    -l.Amount,
			Status:        ItemStatusUnmatched,
			Notes:         "ledger entry missing from payment source",
		}
		items = append(items, item)
		unmatched++
	}

	if err := s.repo.SaveItems(ctx, items); err != nil {
		return nil, err
	}

	run.Status = RunStatusCompleted
	run.TotalTransactions = len(payments) + len(ledgers)
	run.MatchedCount = matched
	run.UnmatchedCount = unmatched
	run.DiscrepancyAmount = discrepancy

	return run, nil
}

func (s *Service) GenerateReport(ctx context.Context, runID string) (*ReconciliationRun, error) {
	return s.repo.GetRun(ctx, runID)
}

func (s *Service) GetUnmatched(ctx context.Context, runID string) ([]*ReconciliationItem, error) {
	items, err := s.repo.GetItemsByRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	var unmatched []*ReconciliationItem
	for _, item := range items {
		if item.Status != ItemStatusMatched {
			unmatched = append(unmatched, item)
		}
	}
	return unmatched, nil
}

func (s *Service) ListRuns(ctx context.Context) ([]*ReconciliationRun, error) {
	return s.repo.ListRuns(ctx)
}
