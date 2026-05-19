package cicd

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name, service string, trigger TriggerType, commitSHA string) (*Pipeline, error) {
	stages := []Stage{
		{Name: "build", Status: StagePending, DurationSeconds: 0, Logs: ""},
		{Name: "test", Status: StagePending, DurationSeconds: 0, Logs: ""},
		{Name: "deploy", Status: StagePending, DurationSeconds: 0, Logs: ""},
	}

	p := &Pipeline{
		ID:        uuid.New().String(),
		Name:      name,
		Service:   service,
		Trigger:   trigger,
		Stages:    stages,
		Status:    StatusPending,
		StartedAt: time.Now(),
		CommitSHA: commitSHA,
	}
	if err := s.repo.Store(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Trigger(ctx context.Context, id string) (*Pipeline, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil || p == nil {
		return nil, err
	}

	p.Status = StatusRunning
	p.StartedAt = time.Now()

	for i := range p.Stages {
		p.Stages[i].Status = StageRunning
		p.Stages[i].Logs = "stage started"
		time.Sleep(5 * time.Millisecond)
		p.Stages[i].Status = StageSuccess
		p.Stages[i].DurationSeconds = 1
		p.Stages[i].Logs = "stage completed"
	}

	p.Status = StatusSuccess
	now := time.Now()
	p.CompletedAt = &now

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetStatus(ctx context.Context, id string) (*Pipeline, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Pipeline, error) {
	return s.repo.List(ctx)
}
