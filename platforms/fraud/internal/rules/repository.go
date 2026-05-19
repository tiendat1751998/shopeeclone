package rules

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, rule *RuleDefinition) error
	Update(ctx context.Context, rule *RuleDefinition) error
	Get(ctx context.Context, id string) (*RuleDefinition, error)
	List(ctx context.Context) ([]*RuleDefinition, error)
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	rules map[string]*RuleDefinition
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		rules: make(map[string]*RuleDefinition),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, rule *RuleDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rule.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	rule.CreatedAt = now
	rule.UpdatedAt = now
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryRepository) Update(ctx context.Context, rule *RuleDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.rules[rule.ID]
	if !ok {
		return ErrRuleNotFound
	}
	rule.CreatedAt = existing.CreatedAt
	rule.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*RuleDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[id]
	if !ok {
		return nil, ErrRuleNotFound
	}
	return rule, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*RuleDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*RuleDefinition, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}
