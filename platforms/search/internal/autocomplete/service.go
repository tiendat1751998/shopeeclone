package autocomplete

import "context"

type Service interface {
	Suggest(ctx context.Context, prefix string, limit int) (*AutocompleteResult, error)
	GetTrending(ctx context.Context, limit int) ([]TrendQuery, error)
	RecordSearch(ctx context.Context, query string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Suggest(ctx context.Context, prefix string, limit int) (*AutocompleteResult, error) {
	suggestions, err := s.repo.Search(ctx, prefix, limit)
	if err != nil {
		return nil, err
	}

	trending, _ := s.repo.GetTrending(ctx, limit)

	seen := make(map[string]bool)
	var combined []Suggestion
	for _, sug := range suggestions {
		if !seen[sug.Text] {
			combined = append(combined, sug)
			seen[sug.Text] = true
		}
	}

	trendingLimit := limit / 3
	for i, t := range trending {
		if i >= trendingLimit {
			break
		}
		if !seen[t.Query] {
			combined = append(combined, Suggestion{
				Text:  t.Query,
				Score: t.Score,
				Type:  "trending",
			})
			seen[t.Query] = true
		}
	}

	if len(combined) > limit {
		combined = combined[:limit]
	}

	return &AutocompleteResult{
		Suggestions: combined,
		TookMs:      0,
	}, nil
}

func (s *service) GetTrending(ctx context.Context, limit int) ([]TrendQuery, error) {
	return s.repo.GetTrending(ctx, limit)
}

func (s *service) RecordSearch(ctx context.Context, query string) error {
	return s.repo.IncrementTrending(ctx, query)
}
