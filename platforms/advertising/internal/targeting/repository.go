package targeting

import (
	"context"
	"sync"
)

type Repository interface {
	GetProfile(ctx context.Context, userID string) (*UserProfile, error)
	StoreProfile(ctx context.Context, profile *UserProfile) error
	GetRules(ctx context.Context, campaignID string) ([]*TargetingRule, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	profiles map[string]*UserProfile
	rules   map[string][]*TargetingRule
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		profiles: make(map[string]*UserProfile),
		rules:    make(map[string][]*TargetingRule),
	}
}

func (r *InMemoryRepository) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.profiles[userID]
	if !ok {
		return &UserProfile{UserID: userID, Devices: []string{"mobile", "desktop"}}, nil
	}
	return p, nil
}

func (r *InMemoryRepository) StoreProfile(ctx context.Context, profile *UserProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profiles[profile.UserID] = profile
	return nil
}

func (r *InMemoryRepository) GetRules(ctx context.Context, campaignID string) ([]*TargetingRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rules[campaignID], nil
}
