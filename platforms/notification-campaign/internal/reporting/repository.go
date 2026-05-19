package reporting

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	UpsertReport(ctx context.Context, r *CampaignReport) error
	GetReport(ctx context.Context, campaignID string) (*CampaignReport, error)
	ListReports(ctx context.Context) ([]*CampaignReport, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	reports map[string]*CampaignReport
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		reports: make(map[string]*CampaignReport),
	}
}

func (r *InMemoryRepository) UpsertReport(ctx context.Context, report *CampaignReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	report.UpdatedAt = time.Now()
	r.reports[report.CampaignID] = report
	return nil
}

func (r *InMemoryRepository) GetReport(ctx context.Context, campaignID string) (*CampaignReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	report, ok := r.reports[campaignID]
	if !ok {
		return nil, nil
	}
	return report, nil
}

func (r *InMemoryRepository) ListReports(ctx context.Context) ([]*CampaignReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*CampaignReport
	for _, report := range r.reports {
		result = append(result, report)
	}
	return result, nil
}
