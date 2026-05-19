package monitoring

import (
	"context"
	"sync"
)

type Repository interface {
	UpsertMetric(ctx context.Context, metric *IndexMetric) error
	GetMetric(ctx context.Context, indexName string) (*IndexMetric, error)
	ListMetrics(ctx context.Context) ([]*IndexMetric, error)
	DeleteMetric(ctx context.Context, indexName string) error
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	metrics map[string]*IndexMetric
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		metrics: make(map[string]*IndexMetric),
	}
}

func (r *InMemoryRepository) UpsertMetric(_ context.Context, metric *IndexMetric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[metric.IndexName] = metric
	return nil
}

func (r *InMemoryRepository) GetMetric(_ context.Context, indexName string) (*IndexMetric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.metrics[indexName]
	if !ok {
		return nil, ErrIndexMetricNotFound
	}
	return m, nil
}

func (r *InMemoryRepository) ListMetrics(_ context.Context) ([]*IndexMetric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*IndexMetric, 0, len(r.metrics))
	for _, m := range r.metrics {
		list = append(list, m)
	}
	return list, nil
}

func (r *InMemoryRepository) DeleteMetric(_ context.Context, indexName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.metrics[indexName]; !ok {
		return ErrIndexMetricNotFound
	}
	delete(r.metrics, indexName)
	return nil
}
