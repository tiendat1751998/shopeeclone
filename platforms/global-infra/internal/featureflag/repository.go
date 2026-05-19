package featureflag

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	Create(ctx context.Context, flag *FeatureFlag) error
	Get(ctx context.Context, name string) (*FeatureFlag, error)
	Update(ctx context.Context, flag *FeatureFlag) error
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]*FeatureFlag, error)
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	flags  map[string]*FeatureFlag
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		flags: make(map[string]*FeatureFlag),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, flag *FeatureFlag) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.flags[flag.Name]; ok {
		return nil
	}
	now := time.Now()
	flag.CreatedAt = now
	flag.UpdatedAt = now
	r.flags[flag.Name] = flag
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, name string) (*FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flag, ok := r.flags[name]
	if !ok {
		return nil, nil
	}
	return flag, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, flag *FeatureFlag) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.flags[flag.Name]; !ok {
		return nil
	}
	flag.UpdatedAt = time.Now()
	r.flags[flag.Name] = flag
	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.flags, name)
	return nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*FeatureFlag, 0, len(r.flags))
	for _, f := range r.flags {
		result = append(result, f)
	}
	return result, nil
}
