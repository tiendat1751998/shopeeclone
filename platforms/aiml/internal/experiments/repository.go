package experiments

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreExperiment(ctx context.Context, exp *Experiment) error
	GetExperiment(ctx context.Context, id string) (*Experiment, error)
	ListExperiments(ctx context.Context) ([]*Experiment, error)
	UpdateExperiment(ctx context.Context, exp *Experiment) error
	StoreResult(ctx context.Context, result *Result) error
	GetResults(ctx context.Context, experimentID string) ([]*Result, error)
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	experiments map[string]*Experiment
	results     map[string][]*Result
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		experiments: make(map[string]*Experiment),
		results:     make(map[string][]*Result),
	}
}

func (r *InMemoryRepository) StoreExperiment(ctx context.Context, exp *Experiment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.experiments[exp.ID]; ok {
		return ErrExperimentExists
	}
	if exp.StartedAt.IsZero() {
		exp.StartedAt = time.Now()
	}
	r.experiments[exp.ID] = exp
	return nil
}

func (r *InMemoryRepository) GetExperiment(ctx context.Context, id string) (*Experiment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	exp, ok := r.experiments[id]
	if !ok {
		return nil, ErrExperimentNotFound
	}
	return exp, nil
}

func (r *InMemoryRepository) ListExperiments(ctx context.Context) ([]*Experiment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Experiment, 0, len(r.experiments))
	for _, exp := range r.experiments {
		result = append(result, exp)
	}
	return result, nil
}

func (r *InMemoryRepository) UpdateExperiment(ctx context.Context, exp *Experiment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.experiments[exp.ID]; !ok {
		return ErrExperimentNotFound
	}
	r.experiments[exp.ID] = exp
	return nil
}

func (r *InMemoryRepository) StoreResult(ctx context.Context, result *Result) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if result.Timestamp.IsZero() {
		result.Timestamp = time.Now()
	}
	r.results[result.ExperimentID] = append(r.results[result.ExperimentID], result)
	return nil
}

func (r *InMemoryRepository) GetResults(ctx context.Context, experimentID string) ([]*Result, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	results := r.results[experimentID]
	if results == nil {
		return []*Result{}, nil
	}
	return results, nil
}
