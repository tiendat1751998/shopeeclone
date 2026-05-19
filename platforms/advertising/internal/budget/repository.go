package budget

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	GetTracker(ctx context.Context, campaignID string) (*BudgetPlan, error)
	UpdateTracker(ctx context.Context, tracker *BudgetPlan) error
	CreateTracker(ctx context.Context, tracker *BudgetPlan) error
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	trackers map[string]*BudgetPlan
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		trackers: make(map[string]*BudgetPlan),
	}
}

func (r *InMemoryRepository) GetTracker(ctx context.Context, campaignID string) (*BudgetPlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.trackers[campaignID]
	if !ok {
		return &BudgetPlan{
			CampaignID:    campaignID,
			DailyBudget:   100,
			LifetimeBudget: 1000,
			SpentToday:    0,
			LifetimeSpent: 0,
			LastResetDate: time.Now().Format("2006-01-02"),
			IsActive:      true,
		}, nil
	}
	return t, nil
}

func (r *InMemoryRepository) UpdateTracker(ctx context.Context, tracker *BudgetPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trackers[tracker.CampaignID] = tracker
	return nil
}

func (r *InMemoryRepository) CreateTracker(ctx context.Context, tracker *BudgetPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trackers[tracker.CampaignID] = tracker
	return nil
}
