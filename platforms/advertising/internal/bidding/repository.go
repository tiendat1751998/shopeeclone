package bidding

import (
	"context"
	"sync"
	"time"
)

type BidHistory struct {
	CampaignID string
	UserID     string
	BidAmount  float64
	Won        bool
	Timestamp  time.Time
}

type Repository interface {
	StoreBidHistory(ctx context.Context, bh *BidHistory) error
	GetBidHistory(ctx context.Context, campaignID string, limit int) ([]*BidHistory, error)
	GetWinRate(ctx context.Context, campaignID string) (float64, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	history []*BidHistory
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		history: make([]*BidHistory, 0),
	}
}

func (r *InMemoryRepository) StoreBidHistory(ctx context.Context, bh *BidHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.history = append(r.history, bh)
	return nil
}

func (r *InMemoryRepository) GetBidHistory(ctx context.Context, campaignID string, limit int) ([]*BidHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*BidHistory
	for i := len(r.history) - 1; i >= 0 && len(result) < limit; i-- {
		if r.history[i].CampaignID == campaignID {
			result = append(result, r.history[i])
		}
	}
	return result, nil
}

func (r *InMemoryRepository) GetWinRate(ctx context.Context, campaignID string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var total, wins int
	for _, h := range r.history {
		if h.CampaignID == campaignID {
			total++
			if h.Won {
				wins++
			}
		}
	}
	if total == 0 {
		return 0, nil
	}
	return float64(wins) / float64(total), nil
}
