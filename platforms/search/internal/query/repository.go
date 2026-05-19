package query

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	LogQuery(ctx context.Context, query string, resultCount int64, tookMs int64) error
	GetRecentQueries(ctx context.Context, limit int) ([]string, error)
	GetFrequentQueries(ctx context.Context, since time.Time, limit int) ([]string, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	logs    []queryLog
}

type queryLog struct {
	query       string
	resultCount int64
	tookMs      int64
	timestamp   time.Time
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

func (r *InMemoryRepository) LogQuery(ctx context.Context, query string, resultCount int64, tookMs int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logs = append(r.logs, queryLog{
		query:       query,
		resultCount: resultCount,
		tookMs:      tookMs,
		timestamp:   time.Now(),
	})

	return nil
}

func (r *InMemoryRepository) GetRecentQueries(ctx context.Context, limit int) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.logs) == 0 {
		return []string{}, nil
	}

	start := len(r.logs) - limit
	if start < 0 {
		start = 0
	}

	var seen = make(map[string]bool)
	var result []string
	for i := len(r.logs) - 1; i >= start; i-- {
		q := r.logs[i].query
		if !seen[q] {
			result = append(result, q)
			seen[q] = true
		}
	}

	return result, nil
}

func (r *InMemoryRepository) GetFrequentQueries(ctx context.Context, since time.Time, limit int) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	freq := make(map[string]int)
	for _, log := range r.logs {
		if log.timestamp.After(since) {
			freq[log.query]++
		}
	}

	type qc struct {
		query string
		count int
	}
	var list []qc
	for q, c := range freq {
		list = append(list, qc{query: q, count: c})
	}

	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j].count > list[i].count {
				list[i], list[j] = list[j], list[i]
			}
		}
	}

	if limit > 0 && len(list) > limit {
		list = list[:limit]
	}

	result := make([]string, len(list))
	for i, item := range list {
		result[i] = item.query
	}

	return result, nil
}
