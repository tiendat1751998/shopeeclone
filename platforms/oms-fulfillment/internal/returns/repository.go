package returns

import (
	"context"
	"sync"
)

type Repository interface {
	Create(ctx context.Context, r *Return) error
	GetByID(ctx context.Context, id string) (*Return, error)
	Update(ctx context.Context, r *Return) error
	List(ctx context.Context) ([]*Return, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	returns map[string]*Return
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{returns: make(map[string]*Return)}
}

func (r *InMemoryRepository) Create(_ context.Context, ret *Return) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.returns[ret.ID] = ret
	return nil
}

func (r *InMemoryRepository) GetByID(_ context.Context, id string) (*Return, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ret, ok := r.returns[id]
	if !ok {
		return nil, ErrReturnNotFound
	}
	return ret, nil
}

func (r *InMemoryRepository) Update(_ context.Context, ret *Return) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.returns[ret.ID]; !ok {
		return ErrReturnNotFound
	}
	r.returns[ret.ID] = ret
	return nil
}

func (r *InMemoryRepository) List(_ context.Context) ([]*Return, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Return
	for _, ret := range r.returns {
		result = append(result, ret)
	}
	return result, nil
}

var _ Repository = (*InMemoryRepository)(nil)
