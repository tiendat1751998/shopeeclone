package analytics

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreImpression(ctx context.Context, imp *Impression) error
	StoreClick(ctx context.Context, click *Click) error
	StoreConversion(ctx context.Context, conv *Conversion) error
	GetImpressions(ctx context.Context, filter *ReportFilter) ([]*Impression, error)
	GetClicks(ctx context.Context, filter *ReportFilter) ([]*Click, error)
	GetConversions(ctx context.Context, filter *ReportFilter) ([]*Conversion, error)
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	impressions []*Impression
	clicks      []*Click
	conversions []*Conversion
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		impressions: make([]*Impression, 0),
		clicks:      make([]*Click, 0),
		conversions: make([]*Conversion, 0),
	}
}

func (r *InMemoryRepository) StoreImpression(ctx context.Context, imp *Impression) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.impressions = append(r.impressions, imp)
	return nil
}

func (r *InMemoryRepository) StoreClick(ctx context.Context, click *Click) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clicks = append(r.clicks, click)
	return nil
}

func (r *InMemoryRepository) StoreConversion(ctx context.Context, conv *Conversion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conversions = append(r.conversions, conv)
	return nil
}

func (r *InMemoryRepository) GetImpressions(ctx context.Context, filter *ReportFilter) ([]*Impression, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Impression
	for _, imp := range r.impressions {
		if filter.CampaignID != "" && imp.CampaignID != filter.CampaignID {
			continue
		}
		if filter.CreativeID != "" && imp.CreativeID != filter.CreativeID {
			continue
		}
		if filter.StartDate != "" {
			start, _ := time.Parse("2006-01-02", filter.StartDate)
			if imp.Timestamp.Before(start) {
				continue
			}
		}
		if filter.EndDate != "" {
			end, _ := time.Parse("2006-01-02", filter.EndDate)
			if imp.Timestamp.After(end.Add(24*time.Hour - time.Second)) {
				continue
			}
		}
		result = append(result, imp)
	}
	return result, nil
}

func (r *InMemoryRepository) GetClicks(ctx context.Context, filter *ReportFilter) ([]*Click, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Click
	for _, cl := range r.clicks {
		if filter.CampaignID != "" && cl.CampaignID != filter.CampaignID {
			continue
		}
		if filter.CreativeID != "" && cl.CreativeID != filter.CreativeID {
			continue
		}
		if filter.StartDate != "" {
			start, _ := time.Parse("2006-01-02", filter.StartDate)
			if cl.Timestamp.Before(start) {
				continue
			}
		}
		if filter.EndDate != "" {
			end, _ := time.Parse("2006-01-02", filter.EndDate)
			if cl.Timestamp.After(end.Add(24*time.Hour - time.Second)) {
				continue
			}
		}
		result = append(result, cl)
	}
	return result, nil
}

func (r *InMemoryRepository) GetConversions(ctx context.Context, filter *ReportFilter) ([]*Conversion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Conversion
	for _, conv := range r.conversions {
		if filter.CampaignID != "" && conv.CampaignID != filter.CampaignID {
			continue
		}
		if filter.CreativeID != "" && conv.CreativeID != filter.CreativeID {
			continue
		}
		if filter.StartDate != "" {
			start, _ := time.Parse("2006-01-02", filter.StartDate)
			if conv.Timestamp.Before(start) {
				continue
			}
		}
		if filter.EndDate != "" {
			end, _ := time.Parse("2006-01-02", filter.EndDate)
			if conv.Timestamp.After(end.Add(24*time.Hour - time.Second)) {
				continue
			}
		}
		result = append(result, conv)
	}
	return result, nil
}
