package search

import "context"

type Service interface {
	Search(ctx context.Context, query SearchQuery) (*SearchResult, error)
	FacetedSearch(ctx context.Context, query SearchQuery) (*SearchResult, error)
	GetByID(ctx context.Context, id string) (*ProductDocument, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Search(ctx context.Context, query SearchQuery) (*SearchResult, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}
	return s.repo.Search(ctx, query)
}

func (s *service) FacetedSearch(ctx context.Context, query SearchQuery) (*SearchResult, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}
	return s.repo.FacetedSearch(ctx, query)
}

func (s *service) GetByID(ctx context.Context, id string) (*ProductDocument, error) {
	return s.repo.GetByID(ctx, id)
}
