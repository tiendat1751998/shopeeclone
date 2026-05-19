package autocomplete

import (
	"context"
	"sort"
	"strings"
	"sync"
)

type Repository interface {
	Store(ctx context.Context, prefix string, suggestion Suggestion) error
	Search(ctx context.Context, prefix string, limit int) ([]Suggestion, error)
	StoreTrending(ctx context.Context, query string, score float64) error
	GetTrending(ctx context.Context, limit int) ([]TrendQuery, error)
	IncrementTrending(ctx context.Context, query string) error
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	prefixes    map[string]map[string]float64
	trending    map[string]float64
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		prefixes: make(map[string]map[string]float64),
		trending: make(map[string]float64),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, prefix string, suggestion Suggestion) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prefix = strings.ToLower(strings.TrimSpace(prefix))
	if prefix == "" {
		return ErrEmptyPrefix
	}

	for i := 1; i <= len(prefix); i++ {
		sub := prefix[:i]
		if r.prefixes[sub] == nil {
			r.prefixes[sub] = make(map[string]float64)
		}
		r.prefixes[sub][strings.ToLower(suggestion.Text)] = suggestion.Score
	}

	return nil
}

func (r *InMemoryRepository) Search(ctx context.Context, prefix string, limit int) ([]Suggestion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prefix = strings.ToLower(strings.TrimSpace(prefix))
	if prefix == "" {
		return nil, ErrEmptyPrefix
	}

	suggestions, ok := r.prefixes[prefix]
	if !ok || len(suggestions) == 0 {
		return []Suggestion{}, nil
	}

	type scored struct {
		text  string
		score float64
	}
	var scoredList []scored
	for text, score := range suggestions {
		scoredList = append(scoredList, scored{text: text, score: score})
	}

	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].score > scoredList[j].score
	})

	if limit <= 0 {
		limit = 10
	}
	if len(scoredList) > limit {
		scoredList = scoredList[:limit]
	}

	result := make([]Suggestion, len(scoredList))
	for i, s := range scoredList {
		result[i] = Suggestion{Text: s.text, Score: s.score, Type: "prefix"}
	}

	return result, nil
}

func (r *InMemoryRepository) StoreTrending(ctx context.Context, query string, score float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.trending[strings.ToLower(strings.TrimSpace(query))] = score
	return nil
}

func (r *InMemoryRepository) GetTrending(ctx context.Context, limit int) ([]TrendQuery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type scored struct {
		query string
		score float64
	}
	var scoredList []scored
	for query, score := range r.trending {
		scoredList = append(scoredList, scored{query: query, score: score})
	}

	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].score > scoredList[j].score
	})

	if limit <= 0 {
		limit = 10
	}
	if len(scoredList) > limit {
		scoredList = scoredList[:limit]
	}

	result := make([]TrendQuery, len(scoredList))
	for i, s := range scoredList {
		result[i] = TrendQuery{Query: s.query, Score: s.score}
	}

	return result, nil
}

func (r *InMemoryRepository) IncrementTrending(ctx context.Context, query string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.trending[strings.ToLower(strings.TrimSpace(query))]++
	return nil
}
