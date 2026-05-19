package behavior

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ProfileRepository interface {
	Save(ctx context.Context, profile *UserBehaviorProfile) error
	Get(ctx context.Context, userID string) (*UserBehaviorProfile, error)
	List(ctx context.Context) ([]*UserBehaviorProfile, error)
}

type RuleRepository interface {
	Create(ctx context.Context, rule *BehavioralRule) error
	List(ctx context.Context) ([]*BehavioralRule, error)
}

type InMemoryProfileRepository struct {
	mu       sync.RWMutex
	profiles map[string]*UserBehaviorProfile
}

func NewInMemoryProfileRepository() *InMemoryProfileRepository {
	return &InMemoryProfileRepository{
		profiles: make(map[string]*UserBehaviorProfile),
	}
}

func (r *InMemoryProfileRepository) Save(ctx context.Context, profile *UserBehaviorProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	profile.LastUpdated = time.Now().UTC()
	r.profiles[profile.UserID] = profile
	return nil
}

func (r *InMemoryProfileRepository) Get(ctx context.Context, userID string) (*UserBehaviorProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	profile, ok := r.profiles[userID]
	if !ok {
		return nil, ErrProfileNotFound
	}
	return profile, nil
}

func (r *InMemoryProfileRepository) List(ctx context.Context) ([]*UserBehaviorProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*UserBehaviorProfile, 0, len(r.profiles))
	for _, p := range r.profiles {
		result = append(result, p)
	}
	return result, nil
}

type InMemoryRuleRepository struct {
	mu    sync.RWMutex
	rules map[string]*BehavioralRule
}

func NewInMemoryRuleRepository() *InMemoryRuleRepository {
	return &InMemoryRuleRepository{
		rules: make(map[string]*BehavioralRule),
	}
}

func (r *InMemoryRuleRepository) Create(ctx context.Context, rule *BehavioralRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rule.ID = uuid.New().String()
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryRuleRepository) List(ctx context.Context) ([]*BehavioralRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*BehavioralRule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}
