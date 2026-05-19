package personalization

import (
	"context"
	"sync"
)

type Repository interface {
	GetProfile(ctx context.Context, userID string) (*UserProfile, error)
	SaveProfile(ctx context.Context, profile *UserProfile) error
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	profiles map[string]*UserProfile
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		profiles: make(map[string]*UserProfile),
	}
}

func (r *InMemoryRepository) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.profiles[userID]
	if !ok {
		return nil, nil
	}
	pCopy := *p
	pCopy.CategoryWeights = make(map[string]float64)
	for k, v := range p.CategoryWeights {
		pCopy.CategoryWeights[k] = v
	}
	pCopy.PreferredBrands = make(map[string]float64)
	for k, v := range p.PreferredBrands {
		pCopy.PreferredBrands[k] = v
	}
	pCopy.InterestVector = make(map[string]float64)
	for k, v := range p.InterestVector {
		pCopy.InterestVector[k] = v
	}
	return &pCopy, nil
}

func (r *InMemoryRepository) SaveProfile(ctx context.Context, profile *UserProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	pCopy := *profile
	pCopy.CategoryWeights = make(map[string]float64)
	for k, v := range profile.CategoryWeights {
		pCopy.CategoryWeights[k] = v
	}
	pCopy.PreferredBrands = make(map[string]float64)
	for k, v := range profile.PreferredBrands {
		pCopy.PreferredBrands[k] = v
	}
	pCopy.InterestVector = make(map[string]float64)
	for k, v := range profile.InterestVector {
		pCopy.InterestVector[k] = v
	}
	r.profiles[profile.UserID] = &pCopy
	return nil
}
