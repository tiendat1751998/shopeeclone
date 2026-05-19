package training

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, job *TrainingJob) error {
	job.Status = StatusPending
	return s.repo.Store(ctx, job)
}

func (s *Service) Get(ctx context.Context, id string) (*TrainingJob, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*TrainingJob, error) {
	return s.repo.List(ctx)
}

func (s *Service) Start(ctx context.Context, id string) error {
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if job.Status != StatusPending {
		return ErrInvalidTransition
	}
	job.Status = StatusRunning
	now := time.Now()
	job.StartedAt = &now
	return s.repo.Update(ctx, job)
}

func (s *Service) Complete(ctx context.Context, id string, metrics map[string]float64) error {
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if job.Status != StatusRunning {
		return ErrInvalidTransition
	}
	job.Status = StatusCompleted
	job.Metrics = metrics
	now := time.Now()
	job.CompletedAt = &now
	return s.repo.Update(ctx, job)
}

func (s *Service) Fail(ctx context.Context, id string, errMsg string) error {
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if job.Status != StatusRunning && job.Status != StatusPending {
		return ErrInvalidTransition
	}
	job.Status = StatusFailed
	job.Error = errMsg
	now := time.Now()
	job.CompletedAt = &now
	return s.repo.Update(ctx, job)
}
