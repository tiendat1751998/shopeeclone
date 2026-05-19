package ledger

import (
	"context"
)

type Service struct {
	accounts AccountRepository
	txns     TransactionRepository
	entries  EntryRepository
}

func NewService(accts AccountRepository, txns TransactionRepository, entries EntryRepository) *Service {
	return &Service{
		accounts: accts,
		txns:     txns,
		entries:  entries,
	}
}

func (s *Service) CreateAccount(ctx context.Context, id, name string, acctType AccountType, code, currency string, initialBalance int64) (*Account, error) {
	a := &Account{
		ID:       id,
		Name:     name,
		Type:     acctType,
		Code:     code,
		Balance:  initialBalance,
		Currency: currency,
	}
	if err := s.accounts.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) PostTransaction(ctx context.Context, debitAccountID, creditAccountID string, amount int64, entryType, referenceType, referenceID, description string) (*Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	debit, err := s.accounts.GetByID(ctx, debitAccountID)
	if err != nil {
		return nil, err
	}
	credit, err := s.accounts.GetByID(ctx, creditAccountID)
	if err != nil {
		return nil, err
	}

	if debit.Balance < amount {
		return nil, ErrInsufficientBalance
	}

	txnID := NewID()
	now := Now()

	debitEntry := &LedgerEntry{
		ID:            NewID(),
		AccountID:     debitAccountID,
		TransactionID: txnID,
		DebitAmount:   amount,
		CreditAmount:  0,
		EntryType:     EntryDebit,
		Description:   description,
		CreatedAt:     now,
	}
	creditEntry := &LedgerEntry{
		ID:            NewID(),
		AccountID:     creditAccountID,
		TransactionID: txnID,
		DebitAmount:   0,
		CreditAmount:  amount,
		EntryType:     EntryCredit,
		Description:   description,
		CreatedAt:     now,
	}

	txn := &Transaction{
		ID:            txnID,
		ReferenceType: referenceType,
		ReferenceID:   referenceID,
		Entries:       []LedgerEntry{*debitEntry, *creditEntry},
		Status:        TxnStatusPosted,
		CreatedAt:     now,
	}

	if err := s.txns.Create(ctx, txn); err != nil {
		return nil, err
	}
	if err := s.entries.Create(ctx, debitEntry); err != nil {
		return nil, err
	}
	if err := s.entries.Create(ctx, creditEntry); err != nil {
		return nil, err
	}

	if err := s.accounts.UpdateBalance(ctx, debitAccountID, debit.Balance-amount); err != nil {
		return nil, err
	}
	if err := s.accounts.UpdateBalance(ctx, creditAccountID, credit.Balance+amount); err != nil {
		return nil, err
	}

	return txn, nil
}

func (s *Service) GetAccountBalance(ctx context.Context, accountID string) (*Account, error) {
	return s.accounts.GetByID(ctx, accountID)
}

func (s *Service) GetAccountStatement(ctx context.Context, accountID string) ([]*LedgerEntry, error) {
	return s.entries.GetByAccount(ctx, accountID)
}

func (s *Service) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.accounts.GetByID(ctx, id)
}

func (s *Service) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	return s.txns.GetByID(ctx, id)
}

func (s *Service) ListTransactions(ctx context.Context, offset, limit int) ([]*Transaction, int64, error) {
	return s.txns.List(ctx, offset, limit)
}

func (s *Service) VoidTransaction(ctx context.Context, id string) (*Transaction, error) {
	txn, err := s.txns.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if txn.Status != TxnStatusPosted {
		return nil, ErrInvalidAmount
	}

	debitEntries, _ := s.entries.GetByTransaction(ctx, id)
	debitTotals := map[string]int64{}
	creditTotals := map[string]int64{}
	for _, e := range debitEntries {
		if e.EntryType == EntryDebit {
			debitTotals[e.AccountID] += e.DebitAmount
		} else {
			creditTotals[e.AccountID] += e.CreditAmount
		}
	}

	for acctID, amt := range debitTotals {
		a, err := s.accounts.GetByID(ctx, acctID)
		if err != nil {
			return nil, err
		}
		if err := s.accounts.UpdateBalance(ctx, acctID, a.Balance+amt); err != nil {
			return nil, err
		}
	}
	for acctID, amt := range creditTotals {
		a, err := s.accounts.GetByID(ctx, acctID)
		if err != nil {
			return nil, err
		}
		if err := s.accounts.UpdateBalance(ctx, acctID, a.Balance-amt); err != nil {
			return nil, err
		}
	}

	if err := s.txns.UpdateStatus(ctx, id, TxnStatusVoided); err != nil {
		return nil, err
	}
	txn.Status = TxnStatusVoided
	return txn, nil
}
