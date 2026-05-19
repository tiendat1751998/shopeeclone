package report_scheduler

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreReport(ctx context.Context, report *ScheduledReport) error
	GetReport(ctx context.Context, id string) (*ScheduledReport, error)
	ListReports(ctx context.Context, organizationID string, offset, limit int) ([]*ScheduledReport, int, error)
	UpdateReport(ctx context.Context, report *ScheduledReport) error
	DeleteReport(ctx context.Context, id string) error
	ListDueReports(ctx context.Context) ([]*ScheduledReport, error)
	StoreGeneration(ctx context.Context, gen *ReportGeneration) error
	GetGeneration(ctx context.Context, id string) (*ReportGeneration, error)
}

type InMemoryRepository struct {
	mu           sync.RWMutex
	reports      map[string]*ScheduledReport
	generations  map[string]*ReportGeneration
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		reports:     make(map[string]*ScheduledReport),
		generations: make(map[string]*ReportGeneration),
	}
}

func (r *InMemoryRepository) StoreReport(ctx context.Context, report *ScheduledReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	if report.CreatedAt.IsZero() {
		report.CreatedAt = now
	}
	report.UpdatedAt = now
	r.reports[report.ID] = report
	return nil
}

func (r *InMemoryRepository) GetReport(ctx context.Context, id string) (*ScheduledReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	report, ok := r.reports[id]
	if !ok {
		return nil, nil
	}
	return report, nil
}

func (r *InMemoryRepository) ListReports(ctx context.Context, organizationID string, offset, limit int) ([]*ScheduledReport, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*ScheduledReport
	for _, rep := range r.reports {
		if organizationID != "" && rep.OrganizationID != organizationID {
			continue
		}
		filtered = append(filtered, rep)
	}
	total := len(filtered)
	if offset >= total {
		return []*ScheduledReport{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return filtered[offset:end], total, nil
}

func (r *InMemoryRepository) UpdateReport(ctx context.Context, report *ScheduledReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	report.UpdatedAt = time.Now()
	r.reports[report.ID] = report
	return nil
}

func (r *InMemoryRepository) DeleteReport(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.reports, id)
	return nil
}

func (r *InMemoryRepository) ListDueReports(ctx context.Context) ([]*ScheduledReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := time.Now()
	var due []*ScheduledReport
	for _, rep := range r.reports {
		if rep.IsActive && rep.NextRunAt.Before(now) {
			due = append(due, rep)
		}
	}
	return due, nil
}

func (r *InMemoryRepository) StoreGeneration(ctx context.Context, gen *ReportGeneration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if gen.GeneratedAt.IsZero() {
		gen.GeneratedAt = time.Now()
	}
	r.generations[gen.ID] = gen
	return nil
}

func (r *InMemoryRepository) GetGeneration(ctx context.Context, id string) (*ReportGeneration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	gen, ok := r.generations[id]
	if !ok {
		return nil, nil
	}
	return gen, nil
}
