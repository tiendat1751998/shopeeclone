package sdk

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

func (s *Service) Register(ctx context.Context, name, language, version, repoURL, docsURL, compatibility string) (*SDK, error) {
	sdk := &SDK{
		ID:               uuid.New().String(),
		Name:             name,
		Language:         language,
		Version:          version,
		RepositoryURL:    repoURL,
		DocumentationURL: docsURL,
		Compatibility:    compatibility,
		IsLatest:         false,
		CreatedAt:        time.Now(),
	}
	if err := s.repo.Store(ctx, sdk); err != nil {
		return nil, err
	}
	return sdk, nil
}

func (s *Service) List(ctx context.Context, language string) ([]*SDK, error) {
	if language != "" {
		return s.repo.ListByLanguage(ctx, language)
	}
	return s.repo.List(ctx)
}

func (s *Service) MarkLatest(ctx context.Context, id string) (*SDK, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil || target == nil {
		return nil, err
	}

	all, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, sdk := range all {
		if sdk.Language == target.Language && sdk.IsLatest {
			sdk.IsLatest = false
			if err := s.repo.Update(ctx, sdk); err != nil {
				return nil, err
			}
		}
	}

	target.IsLatest = true
	if err := s.repo.Update(ctx, target); err != nil {
		return nil, err
	}
	return target, nil
}
