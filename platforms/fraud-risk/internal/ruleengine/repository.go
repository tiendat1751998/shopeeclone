package ruleengine

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RuleRepository interface {
	Create(ctx context.Context, rule *Rule) error
	Get(ctx context.Context, id string) (*Rule, error)
	List(ctx context.Context) ([]*Rule, error)
	Update(ctx context.Context, rule *Rule) error
}

type RuleSetRepository interface {
	Create(ctx context.Context, rs *RuleSet) error
	Get(ctx context.Context, id string) (*RuleSet, error)
	List(ctx context.Context) ([]*RuleSet, error)
}

type InMemoryRuleRepository struct {
	mu    sync.RWMutex
	rules map[string]*Rule
}

func NewInMemoryRuleRepository() *InMemoryRuleRepository {
	return &InMemoryRuleRepository{
		rules: make(map[string]*Rule),
	}
}

func (r *InMemoryRuleRepository) Create(ctx context.Context, rule *Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rule.ID = uuid.New().String()
	now := time.Now().UTC()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryRuleRepository) Get(ctx context.Context, id string) (*Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[id]
	if !ok {
		return nil, ErrRuleNotFound
	}
	return rule, nil
}

func (r *InMemoryRuleRepository) List(ctx context.Context) ([]*Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}

func (r *InMemoryRuleRepository) Update(ctx context.Context, rule *Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.rules[rule.ID]
	if !ok {
		return ErrRuleNotFound
	}
	rule.CreatedAt = existing.CreatedAt
	rule.UpdatedAt = time.Now().UTC()
	r.rules[rule.ID] = rule
	return nil
}

type InMemoryRuleSetRepository struct {
	mu       sync.RWMutex
	rulesets map[string]*RuleSet
}

func NewInMemoryRuleSetRepository() *InMemoryRuleSetRepository {
	return &InMemoryRuleSetRepository{
		rulesets: make(map[string]*RuleSet),
	}
}

func (r *InMemoryRuleSetRepository) Create(ctx context.Context, rs *RuleSet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rs.ID = uuid.New().String()
	now := time.Now().UTC()
	rs.CreatedAt = now
	rs.UpdatedAt = now
	r.rulesets[rs.ID] = rs
	return nil
}

func (r *InMemoryRuleSetRepository) Get(ctx context.Context, id string) (*RuleSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rs, ok := r.rulesets[id]
	if !ok {
		return nil, ErrRuleSetNotFound
	}
	return rs, nil
}

func (r *InMemoryRuleSetRepository) List(ctx context.Context) ([]*RuleSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*RuleSet, 0, len(r.rulesets))
	for _, rs := range r.rulesets {
		result = append(result, rs)
	}
	return result, nil
}
