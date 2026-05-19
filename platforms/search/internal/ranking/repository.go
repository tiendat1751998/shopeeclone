package ranking

import (
	"context"
	"sync"
)

type Repository interface {
	GetConfig(ctx context.Context) (*RankingConfig, error)
	SetConfig(ctx context.Context, config *RankingConfig) error
	RecordClick(ctx context.Context, productID, query string) error
	GetClickSignals(ctx context.Context, productID string) (*ClickSignal, error)
	GetTopClicked(ctx context.Context, limit int) ([]ClickSignal, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	config  *RankingConfig
	clicks  map[string]*ClickSignal
}

func NewInMemoryRepository() *InMemoryRepository {
	cfg := DefaultRankingConfig()
	return &InMemoryRepository{
		config: &cfg,
		clicks: make(map[string]*ClickSignal),
	}
}

func (r *InMemoryRepository) GetConfig(ctx context.Context) (*RankingConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cfg := *r.config
	return &cfg, nil
}

func (r *InMemoryRepository) SetConfig(ctx context.Context, config *RankingConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.config = config
	return nil
}

func (r *InMemoryRepository) RecordClick(ctx context.Context, productID, query string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := productID + ":" + query
	signal, ok := r.clicks[key]
	if !ok {
		signal = &ClickSignal{
			ProductID: productID,
			Query:     query,
		}
		r.clicks[key] = signal
	}

	signal.Count++
	totalClicks := int64(0)
	for _, s := range r.clicks {
		if s.ProductID == productID {
			totalClicks += s.Count
		}
	}

	if totalClicks > 0 {
		signal.CTR = float64(signal.Count) / float64(totalClicks)
	}

	return nil
}

func (r *InMemoryRepository) GetClickSignals(ctx context.Context, productID string) (*ClickSignal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var combined ClickSignal
	combined.ProductID = productID
	totalCount := int64(0)

	for _, signal := range r.clicks {
		if signal.ProductID == productID {
			combined.Count += signal.Count
			totalCount += signal.Count
		}
	}

	if totalCount > 0 {
		combined.CTR = float64(combined.Count) / float64(totalCount)
	}

	return &combined, nil
}

func (r *InMemoryRepository) GetTopClicked(ctx context.Context, limit int) ([]ClickSignal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	productClicks := make(map[string]int64)
	productQueries := make(map[string]string)

	for _, signal := range r.clicks {
		productClicks[signal.ProductID] += signal.Count
		productQueries[signal.ProductID] = signal.Query
	}

	type pc struct {
		productID string
		count     int64
		query     string
	}
	var list []pc
	for pid, count := range productClicks {
		list = append(list, pc{productID: pid, count: count, query: productQueries[pid]})
	}

	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j].count > list[i].count {
				list[i], list[j] = list[j], list[i]
			}
		}
	}

	if limit <= 0 || limit > len(list) {
		limit = len(list)
	}
	list = list[:limit]

	result := make([]ClickSignal, len(list))
	for i, item := range list {
		result[i] = ClickSignal{
			ProductID: item.productID,
			Query:     item.query,
			Count:     item.count,
		}
	}

	return result, nil
}
