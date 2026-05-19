package budget

import (
	"context"
	"time"
)

type Service interface {
	CheckBudget(ctx context.Context, campaignID string, bidAmount float64) (bool, error)
	DeductSpend(ctx context.Context, campaignID string, amount float64) error
	GetRemainingBudget(ctx context.Context, campaignID string) (daily float64, lifetime float64, err error)
	ResetDailyBudget(ctx context.Context, campaignID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CheckBudget(ctx context.Context, campaignID string, bidAmount float64) (bool, error) {
	tracker, err := s.repo.GetTracker(ctx, campaignID)
	if err != nil {
		return false, err
	}
	if !tracker.IsActive {
		return false, nil
	}

	today := time.Now().Format("2006-01-02")
	if tracker.LastResetDate != today {
		tracker.SpentToday = 0
		tracker.LastResetDate = today
		s.repo.UpdateTracker(ctx, tracker)
	}

	if tracker.DailyBudget > 0 && tracker.SpentToday+bidAmount > tracker.DailyBudget {
		return false, nil
	}
	if tracker.LifetimeBudget > 0 && tracker.LifetimeSpent+bidAmount > tracker.LifetimeBudget {
		return false, nil
	}
	return true, nil
}

func (s *service) DeductSpend(ctx context.Context, campaignID string, amount float64) error {
	tracker, err := s.repo.GetTracker(ctx, campaignID)
	if err != nil {
		return err
	}
	tracker.SpentToday += amount
	tracker.LifetimeSpent += amount
	return s.repo.UpdateTracker(ctx, tracker)
}

func (s *service) GetRemainingBudget(ctx context.Context, campaignID string) (float64, float64, error) {
	tracker, err := s.repo.GetTracker(ctx, campaignID)
	if err != nil {
		return 0, 0, err
	}
	dailyRemaining := tracker.DailyBudget - tracker.SpentToday
	lifetimeRemaining := tracker.LifetimeBudget - tracker.LifetimeSpent
	if dailyRemaining < 0 {
		dailyRemaining = 0
	}
	if lifetimeRemaining < 0 {
		lifetimeRemaining = 0
	}
	return dailyRemaining, lifetimeRemaining, nil
}

func (s *service) ResetDailyBudget(ctx context.Context, campaignID string) error {
	tracker, err := s.repo.GetTracker(ctx, campaignID)
	if err != nil {
		return err
	}
	tracker.SpentToday = 0
	tracker.LastResetDate = time.Now().Format("2006-01-02")
	return s.repo.UpdateTracker(ctx, tracker)
}
