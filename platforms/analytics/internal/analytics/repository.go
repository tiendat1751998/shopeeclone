package analytics

import (
	"context"
	"sync"
)

type Repository interface {
	StoreQueryResult(ctx context.Context, report *Report) error
	GetQueryResult(ctx context.Context, reportID string) (*Report, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	reports map[string]*Report
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		reports: make(map[string]*Report),
	}
}

func (r *InMemoryRepository) StoreQueryResult(ctx context.Context, report *Report) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reports[report.ID] = report
	return nil
}

func (r *InMemoryRepository) GetQueryResult(ctx context.Context, reportID string) (*Report, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	report, ok := r.reports[reportID]
	if !ok {
		return nil, ErrMetricNotFound
	}
	return report, nil
}
