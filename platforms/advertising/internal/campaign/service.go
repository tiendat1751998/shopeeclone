package campaign

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, c *Campaign) (*Campaign, error)
	Get(ctx context.Context, id string) (*Campaign, error)
	List(ctx context.Context, status CampaignStatus) ([]*Campaign, error)
	Update(ctx context.Context, c *Campaign) (*Campaign, error)
	Pause(ctx context.Context, id string) (*Campaign, error)
	Resume(ctx context.Context, id string) (*Campaign, error)
	End(ctx context.Context, id string) (*Campaign, error)
	GetPerformance(ctx context.Context, id string) (*Performance, error)
}

type Performance struct {
	CampaignID  string
	Impressions int64
	Clicks      int64
	Conversions int64
	Spend       float64
	Revenue     float64
	CTR         float64
	CVR         float64
	CPC         float64
	CPM         float64
	ROAS        float64
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, c *Campaign) (*Campaign, error) {
	if c.Name == "" {
		return nil, ErrEmptyName
	}
	if c.Type != CampaignTypeCPC && c.Type != CampaignTypeCPM && c.Type != CampaignTypeCPA {
		return nil, ErrInvalidType
	}
	if c.Budget.Daily < 0 || c.Budget.Lifetime < 0 {
		return nil, ErrInvalidBudget
	}
	if !c.DateRange.End.IsZero() && c.DateRange.End.Before(c.DateRange.Start) {
		return nil, ErrInvalidDateRange
	}
	c.ID = uuid.New().String()
	c.Status = CampaignStatusDraft
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	if c.QualityScore == 0 {
		c.QualityScore = 1.0
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) Get(ctx context.Context, id string) (*Campaign, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) List(ctx context.Context, status CampaignStatus) ([]*Campaign, error) {
	return s.repo.List(ctx, status)
}

func (s *service) Update(ctx context.Context, c *Campaign) (*Campaign, error) {
	existing, err := s.repo.Get(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	if existing.Status == CampaignStatusEnded {
		return nil, ErrInvalidStatus
	}
	c.UpdatedAt = time.Now()
	return s.repo.Update(ctx, c)
}

func (s *service) Pause(ctx context.Context, id string) (*Campaign, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.Status != CampaignStatusActive {
		return nil, ErrInvalidStatus
	}
	c.Status = CampaignStatusPaused
	c.UpdatedAt = time.Now()
	return s.repo.Update(ctx, c)
}

func (s *service) Resume(ctx context.Context, id string) (*Campaign, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.Status != CampaignStatusPaused {
		return nil, ErrInvalidStatus
	}
	c.Status = CampaignStatusActive
	c.UpdatedAt = time.Now()
	return s.repo.Update(ctx, c)
}

func (s *service) End(ctx context.Context, id string) (*Campaign, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.Status == CampaignStatusEnded {
		return nil, ErrInvalidStatus
	}
	c.Status = CampaignStatusEnded
	c.UpdatedAt = time.Now()
	return s.repo.Update(ctx, c)
}

func (s *service) GetPerformance(ctx context.Context, id string) (*Performance, error) {
	return s.repo.GetPerformance(ctx, id)
}
