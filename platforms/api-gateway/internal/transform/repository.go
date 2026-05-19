package transform

import "sync"

type Repository interface {
	Store(rule *Rule) error
	Get(id string) (*Rule, error)
	List() ([]*Rule, error)
	Delete(id string) error
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	rules map[string]*Rule
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		rules: make(map[string]*Rule),
	}
}

func (r *InMemoryRepository) Store(rule *Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryRepository) Get(id string) (*Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[id]
	if !ok {
		return nil, nil
	}
	return rule, nil
}

func (r *InMemoryRepository) List() ([]*Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}

func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, id)
	return nil
}
