package ratelimit

import (
	"context"
	"sync"
)

type Repository interface {
	StoreRule(ctx context.Context, rule *RateLimitRule) error
	GetRule(ctx context.Context, keyPattern string) (*RateLimitRule, error)
	ListRules(ctx context.Context) ([]*RateLimitRule, error)
	DeleteRule(ctx context.Context, keyPattern string) error
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	rules map[string]*RateLimitRule
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		rules: make(map[string]*RateLimitRule),
	}
}

func (r *InMemoryRepository) StoreRule(ctx context.Context, rule *RateLimitRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.KeyPattern] = rule
	return nil
}

func (r *InMemoryRepository) GetRule(ctx context.Context, keyPattern string) (*RateLimitRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[keyPattern]
	if !ok {
		return nil, nil
	}
	return rule, nil
}

func (r *InMemoryRepository) ListRules(ctx context.Context) ([]*RateLimitRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*RateLimitRule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}

func (r *InMemoryRepository) DeleteRule(ctx context.Context, keyPattern string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, keyPattern)
	return nil
}
