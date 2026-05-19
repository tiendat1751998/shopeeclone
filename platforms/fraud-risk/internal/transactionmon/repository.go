package transactionmon

import (
	"context"
	"sync"
)

type Repository interface {
	Save(ctx context.Context, m *TransactionMonitor) error
	Get(ctx context.Context, userID string) (*TransactionMonitor, error)
	List(ctx context.Context) ([]*TransactionMonitor, error)
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	monitors  map[string]*TransactionMonitor
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		monitors: make(map[string]*TransactionMonitor),
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, m *TransactionMonitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.monitors[m.UserID] = m
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, userID string) (*TransactionMonitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.monitors[userID]
	if !ok {
		return nil, ErrMonitorNotFound
	}
	return m, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*TransactionMonitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*TransactionMonitor, 0, len(r.monitors))
	for _, m := range r.monitors {
		result = append(result, m)
	}
	return result, nil
}
