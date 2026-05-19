package creative

import (
	"context"
	"math/rand"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, c *Creative) (*Creative, error)
	Get(ctx context.Context, id string) (*Creative, error)
	GetByCampaign(ctx context.Context, campaignID string) ([]*Creative, error)
	List(ctx context.Context, status CreativeStatus) ([]*Creative, error)
	Approve(ctx context.Context, id string) (*Creative, error)
	Reject(ctx context.Context, id string) (*Creative, error)
	ServeCreative(ctx context.Context, campaignID string) (*Creative, error)
	RotateCreative(ctx context.Context, creatives []*Creative) (*Creative, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, c *Creative) (*Creative, error) {
	c.ID = uuid.New().String()
	c.Status = CreativeStatusDraft
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return s.repo.Create(ctx, c)
}

func (s *service) Get(ctx context.Context, id string) (*Creative, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetByCampaign(ctx context.Context, campaignID string) ([]*Creative, error) {
	return s.repo.GetByCampaign(ctx, campaignID)
}

func (s *service) List(ctx context.Context, status CreativeStatus) ([]*Creative, error) {
	return s.repo.List(ctx, status)
}

func (s *service) Approve(ctx context.Context, id string) (*Creative, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Status = CreativeStatusApproved
	c.UpdatedAt = time.Now()
	return s.repo.Update(ctx, c)
}

func (s *service) Reject(ctx context.Context, id string) (*Creative, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Status = CreativeStatusRejected
	c.UpdatedAt = time.Now()
	return s.repo.Update(ctx, c)
}

func (s *service) ServeCreative(ctx context.Context, campaignID string) (*Creative, error) {
	creatives, err := s.repo.GetByCampaign(ctx, campaignID)
	if err != nil || len(creatives) == 0 {
		return nil, err
	}

	var approved []*Creative
	for _, cr := range creatives {
		if cr.Status == CreativeStatusApproved {
			approved = append(approved, cr)
		}
	}
	if len(approved) == 0 {
		return nil, nil
	}

	return s.RotateCreative(ctx, approved)
}

func (s *service) RotateCreative(ctx context.Context, creatives []*Creative) (*Creative, error) {
	if len(creatives) == 0 {
		return nil, nil
	}
	if len(creatives) == 1 {
		return creatives[0], nil
	}

	sort.Slice(creatives, func(i, j int) bool {
		return creatives[i].Performance.CTR > creatives[j].Performance.CTR
	})

	if rand.Float64() < 0.3 {
		return creatives[rand.Intn(len(creatives))], nil
	}
	return creatives[0], nil
}
