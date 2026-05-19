package decisionlog

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, log *DecisionLog) error
	Get(ctx context.Context, id string) (*DecisionLog, error)
	List(ctx context.Context) ([]*DecisionLog, error)
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	logs  map[string]*DecisionLog
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		logs: make(map[string]*DecisionLog),
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, log *DecisionLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	log.ID = uuid.New().String()
	log.Timestamp = time.Now().UTC()
	r.logs[log.ID] = log
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*DecisionLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	log, ok := r.logs[id]
	if !ok {
		return nil, ErrDecisionNotFound
	}
	return log, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*DecisionLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*DecisionLog, 0, len(r.logs))
	for _, log := range r.logs {
		result = append(result, log)
	}
	return result, nil
}
