package traffic

import (
	"context"
	"sync"
)

type Repository interface {
	CreateRule(ctx context.Context, rule *TrafficRule) error
	ListRules(ctx context.Context) ([]*TrafficRule, error)
	GetRule(ctx context.Context, id string) (*TrafficRule, error)
	DeleteRule(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	rules map[string]*TrafficRule
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		rules: make(map[string]*TrafficRule),
	}
}

func (r *InMemoryRepository) CreateRule(ctx context.Context, rule *TrafficRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryRepository) ListRules(ctx context.Context) ([]*TrafficRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*TrafficRule
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}

func (r *InMemoryRepository) GetRule(ctx context.Context, id string) (*TrafficRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rule, ok := r.rules[id]
	if !ok {
		return nil, nil
	}
	return rule, nil
}

func (r *InMemoryRepository) DeleteRule(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, id)
	return nil
}
