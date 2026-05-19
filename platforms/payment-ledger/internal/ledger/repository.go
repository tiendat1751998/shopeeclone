package ledger

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, a *Account) error
	GetByID(ctx context.Context, id string) (*Account, error)
	UpdateBalance(ctx context.Context, id string, balance int64) error
	List(ctx context.Context) ([]*Account, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, t *Transaction) error
	GetByID(ctx context.Context, id string) (*Transaction, error)
	UpdateStatus(ctx context.Context, id string, status TransactionStatus) error
	List(ctx context.Context, offset, limit int) ([]*Transaction, int64, error)
}

type EntryRepository interface {
	Create(ctx context.Context, e *LedgerEntry) error
	GetByTransaction(ctx context.Context, txnID string) ([]*LedgerEntry, error)
	GetByAccount(ctx context.Context, accountID string) ([]*LedgerEntry, error)
}

type InMemoryAccountRepo struct {
	mu   sync.RWMutex
	data map[string]*Account
}

func NewInMemoryAccountRepo() *InMemoryAccountRepo {
	return &InMemoryAccountRepo{data: make(map[string]*Account)}
}

func (r *InMemoryAccountRepo) Create(ctx context.Context, a *Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[a.ID]; ok {
		return ErrAccountAlreadyExists
	}
	r.data[a.ID] = a
	return nil
}

func (r *InMemoryAccountRepo) GetByID(ctx context.Context, id string) (*Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.data[id]
	if !ok {
		return nil, ErrAccountNotFound
	}
	return a, nil
}

func (r *InMemoryAccountRepo) UpdateBalance(ctx context.Context, id string, balance int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.data[id]
	if !ok {
		return ErrAccountNotFound
	}
	a.Balance = balance
	return nil
}

func (r *InMemoryAccountRepo) List(ctx context.Context) ([]*Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Account, 0, len(r.data))
	for _, a := range r.data {
		result = append(result, a)
	}
	return result, nil
}

type InMemoryTransactionRepo struct {
	mu   sync.RWMutex
	data map[string]*Transaction
}

func NewInMemoryTransactionRepo() *InMemoryTransactionRepo {
	return &InMemoryTransactionRepo{data: make(map[string]*Transaction)}
}

func (r *InMemoryTransactionRepo) Create(ctx context.Context, t *Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[t.ID] = t
	return nil
}

func (r *InMemoryTransactionRepo) GetByID(ctx context.Context, id string) (*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok {
		return nil, ErrTransactionNotFound
	}
	return t, nil
}

func (r *InMemoryTransactionRepo) UpdateStatus(ctx context.Context, id string, status TransactionStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.data[id]
	if !ok {
		return ErrTransactionNotFound
	}
	t.Status = status
	return nil
}

func (r *InMemoryTransactionRepo) List(ctx context.Context, offset, limit int) ([]*Transaction, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := int64(len(r.data))
	items := make([]*Transaction, 0, limit)
	i := 0
	for _, t := range r.data {
		if i >= offset && len(items) < limit {
			items = append(items, t)
		}
		i++
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, total, nil
}

type InMemoryEntryRepo struct {
	mu   sync.RWMutex
	data map[string]*LedgerEntry
}

func NewInMemoryEntryRepo() *InMemoryEntryRepo {
	return &InMemoryEntryRepo{data: make(map[string]*LedgerEntry)}
}

func (r *InMemoryEntryRepo) Create(ctx context.Context, e *LedgerEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[e.ID] = e
	return nil
}

func (r *InMemoryEntryRepo) GetByTransaction(ctx context.Context, txnID string) ([]*LedgerEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*LedgerEntry, 0)
	for _, e := range r.data {
		if e.TransactionID == txnID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *InMemoryEntryRepo) GetByAccount(ctx context.Context, accountID string) ([]*LedgerEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*LedgerEntry, 0)
	for _, e := range r.data {
		if e.AccountID == accountID {
			result = append(result, e)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt < result[j].CreatedAt
	})
	return result, nil
}

func NewID() string {
	return uuid.New().String()
}

func Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
