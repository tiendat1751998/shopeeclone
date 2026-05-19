package ratelimit

import "sync"

type Repository interface {
	StoreRule(rule *RateLimitRule) error
	GetRule(key string) (*RateLimitRule, error)
	ListRules() ([]*RateLimitRule, error)
	DeleteRule(key string) error
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

func (r *InMemoryRepository) StoreRule(rule *RateLimitRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.Key] = rule
	return nil
}

func (r *InMemoryRepository) GetRule(key string) (*RateLimitRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[key]
	if !ok {
		return nil, nil
	}
	return rule, nil
}

func (r *InMemoryRepository) ListRules() ([]*RateLimitRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*RateLimitRule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}

func (r *InMemoryRepository) DeleteRule(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, key)
	return nil
}
