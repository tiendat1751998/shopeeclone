package campaign

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("campaign not found")
	ErrInvalidStatus = errors.New("invalid status transition")
)

type Service interface {
	Create(ctx context.Context, req *CreateCampaignRequest) (*Campaign, error)
	Get(ctx context.Context, id string) (*Campaign, error)
	List(ctx context.Context) ([]*Campaign, error)
	Update(ctx context.Context, id string, req *UpdateCampaignRequest) (*Campaign, error)
	Start(ctx context.Context, id string) error
	Pause(ctx context.Context, id string) error
	Resume(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string) error
	ExecuteSchedule(ctx context.Context) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *CreateCampaignRequest) (*Campaign, error) {
	c := &Campaign{
		Name:            req.Name,
		Type:            req.Type,
		Channel:         req.Channel,
		Schedule:        req.Schedule,
		AudienceQuery:   req.AudienceQuery,
		ContentTemplate: req.ContentTemplate,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) Get(ctx context.Context, id string) (*Campaign, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *service) List(ctx context.Context) ([]*Campaign, error) {
	return s.repo.List(ctx)
}

func (s *service) Update(ctx context.Context, id string, req *UpdateCampaignRequest) (*Campaign, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Type != nil {
		c.Type = *req.Type
	}
	if req.Channel != nil {
		c.Channel = *req.Channel
	}
	if req.Schedule != nil {
		c.Schedule = *req.Schedule
	}
	if req.AudienceQuery != nil {
		c.AudienceQuery = *req.AudienceQuery
	}
	if req.ContentTemplate != nil {
		c.ContentTemplate = *req.ContentTemplate
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) transition(ctx context.Context, c *Campaign, newStatus Status) error {
	validNext, ok := ValidTransitions[c.Status]
	if !ok {
		return ErrInvalidStatus
	}
	found := false
	for _, s := range validNext {
		if s == newStatus {
			found = true
			break
		}
	}
	if !found {
		return ErrInvalidStatus
	}
	c.Status = newStatus
	return s.repo.Update(ctx, c)
}

func (s *service) Start(ctx context.Context, id string) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrNotFound
	}
	return s.transition(ctx, c, StatusRunning)
}

func (s *service) Pause(ctx context.Context, id string) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrNotFound
	}
	return s.transition(ctx, c, StatusPaused)
}

func (s *service) Resume(ctx context.Context, id string) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrNotFound
	}
	return s.transition(ctx, c, StatusRunning)
}

func (s *service) Cancel(ctx context.Context, id string) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrNotFound
	}
	return s.transition(ctx, c, StatusCancelled)
}

func (s *service) ExecuteSchedule(ctx context.Context) error {
	campaigns, err := s.repo.List(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	for _, c := range campaigns {
		if c.Status == StatusScheduled && now.After(c.Schedule.StartAt) {
			c.Status = StatusRunning
			if err := s.repo.Update(ctx, c); err != nil {
				return err
			}
		}
		if c.Status == StatusRunning && !c.Schedule.EndAt.IsZero() && now.After(c.Schedule.EndAt) {
			c.Status = StatusCompleted
			if err := s.repo.Update(ctx, c); err != nil {
				return err
			}
		}
	}
	return nil
}
