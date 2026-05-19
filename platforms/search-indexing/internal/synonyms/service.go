package synonyms

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	CreateSet(ctx context.Context, words []string, language string) (*SynonymSet, error)
	ExpandQuery(ctx context.Context, query string) ([]string, error)
	GetSynonyms(ctx context.Context, word string) ([]string, error)
	GetSet(ctx context.Context, id string) (*SynonymSet, error)
	ListSets(ctx context.Context) ([]*SynonymSet, error)
	RemoveSet(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateSet(ctx context.Context, words []string, language string) (*SynonymSet, error) {
	if len(words) == 0 {
		return nil, ErrEmptyWords
	}
	set := &SynonymSet{
		ID:       uuid.New().String(),
		Words:    words,
		Language: language,
		IsActive: true,
	}
	if err := s.repo.CreateSet(ctx, set); err != nil {
		return nil, err
	}
	return set, nil
}

func (s *service) ExpandQuery(ctx context.Context, query string) ([]string, error) {
	graph, err := s.repo.GetGraph(ctx)
	if err != nil {
		return nil, err
	}

	words := strings.Fields(strings.ToLower(query))
	expanded := make(map[string]bool)

	for _, word := range words {
		expanded[word] = true
		if syns, ok := graph.Edges[word]; ok {
			for _, syn := range syns {
				expanded[syn] = true
			}
		}
	}

	result := make([]string, 0, len(expanded))
	for w := range expanded {
		result = append(result, w)
	}
	return result, nil
}

func (s *service) GetSynonyms(ctx context.Context, word string) ([]string, error) {
	graph, err := s.repo.GetGraph(ctx)
	if err != nil {
		return nil, err
	}
	syns := graph.Edges[strings.ToLower(word)]
	if syns == nil {
		return []string{}, nil
	}
	return syns, nil
}

func (s *service) GetSet(ctx context.Context, id string) (*SynonymSet, error) {
	return s.repo.GetSet(ctx, id)
}

func (s *service) ListSets(ctx context.Context) ([]*SynonymSet, error) {
	return s.repo.ListSets(ctx)
}

func (s *service) RemoveSet(ctx context.Context, id string) error {
	return s.repo.DeleteSet(ctx, id)
}
