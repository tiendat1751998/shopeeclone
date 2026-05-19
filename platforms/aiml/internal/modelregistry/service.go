package modelregistry

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, model *Model) error {
	return s.repo.Store(ctx, model)
}

func (s *Service) Get(ctx context.Context, id string) (*Model, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Model, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListByStage(ctx context.Context, stage Stage) ([]*Model, error) {
	return s.repo.ListByStage(ctx, stage)
}

func (s *Service) Promote(ctx context.Context, id string, stage Stage) error {
	model, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if model.Status == StageArchived {
		return ErrInvalidStage
	}
	model.Status = stage
	return s.repo.Update(ctx, model)
}

func (s *Service) Archive(ctx context.Context, id string) error {
	model, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	model.Status = StageArchived
	return s.repo.Update(ctx, model)
}
