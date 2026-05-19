package docs

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

func (s *Service) Create(ctx context.Context, title, content, service, category string, tags []string, version string) (*DocPage, error) {
	now := time.Now()
	doc := &DocPage{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		Service:   service,
		Category:  category,
		Tags:      tags,
		Version:   version,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Store(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*DocPage, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, service, category string) ([]*DocPage, error) {
	if service != "" {
		return s.repo.ListByService(ctx, service)
	}
	if category != "" {
		return s.repo.ListByCategory(ctx, category)
	}
	return s.repo.List(ctx)
}

func (s *Service) Search(ctx context.Context, query string) ([]*DocPage, error) {
	return s.repo.Search(ctx, query)
}

func (s *Service) Update(ctx context.Context, id, title, content, service, category string, tags []string, version string) (*DocPage, error) {
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}
	if title != "" {
		doc.Title = title
	}
	if content != "" {
		doc.Content = content
	}
	if service != "" {
		doc.Service = service
	}
	if category != "" {
		doc.Category = category
	}
	if tags != nil {
		doc.Tags = tags
	}
	if version != "" {
		doc.Version = version
	}
	doc.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}
