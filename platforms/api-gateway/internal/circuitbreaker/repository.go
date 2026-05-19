package circuitbreaker

import "sync"

type Repository interface {
	Store(cb *CircuitBreaker) error
	Get(id string) (*CircuitBreaker, error)
	List() ([]*CircuitBreaker, error)
	Delete(id string) error
}

type InMemoryRepository struct {
	mu  sync.RWMutex
	cbs map[string]*CircuitBreaker
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		cbs: make(map[string]*CircuitBreaker),
	}
}

func (r *InMemoryRepository) Store(cb *CircuitBreaker) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cbs[cb.ID] = cb
	return nil
}

func (r *InMemoryRepository) Get(id string) (*CircuitBreaker, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cb, ok := r.cbs[id]
	if !ok {
		return nil, nil
	}
	return cb, nil
}

func (r *InMemoryRepository) List() ([]*CircuitBreaker, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*CircuitBreaker, 0, len(r.cbs))
	for _, cb := range r.cbs {
		result = append(result, cb)
	}
	return result, nil
}

func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cbs, id)
	return nil
}
