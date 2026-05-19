package reconciliation

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateRun(ctx context.Context, r *ReconciliationRun) error
	GetRun(ctx context.Context, id string) (*ReconciliationRun, error)
	ListRuns(ctx context.Context) ([]*ReconciliationRun, error)
	SaveItems(ctx context.Context, items []*ReconciliationItem) error
	GetItemsByRun(ctx context.Context, runID string) ([]*ReconciliationItem, error)
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	runs  map[string]*ReconciliationRun
	items map[string][]*ReconciliationItem
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		runs:  make(map[string]*ReconciliationRun),
		items: make(map[string][]*ReconciliationItem),
	}
}

func (r *InMemoryRepository) CreateRun(ctx context.Context, run *ReconciliationRun) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if run.ID == "" {
		run.ID = uuid.New().String()
	}
	r.runs[run.ID] = run
	return nil
}

func (r *InMemoryRepository) GetRun(ctx context.Context, id string) (*ReconciliationRun, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	run, ok := r.runs[id]
	if !ok {
		return nil, ErrRunNotFound
	}
	return run, nil
}

func (r *InMemoryRepository) ListRuns(ctx context.Context) ([]*ReconciliationRun, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*ReconciliationRun, 0, len(r.runs))
	for _, run := range r.runs {
		result = append(result, run)
	}
	return result, nil
}

func (r *InMemoryRepository) SaveItems(ctx context.Context, items []*ReconciliationItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range items {
		r.items[item.RunID] = append(r.items[item.RunID], item)
	}
	return nil
}

func (r *InMemoryRepository) GetItemsByRun(ctx context.Context, runID string) ([]*ReconciliationItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items, ok := r.items[runID]
	if !ok {
		return nil, nil
	}
	return items, nil
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
