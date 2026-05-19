package dashboard

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreDashboard(ctx context.Context, d *Dashboard) error
	GetDashboard(ctx context.Context, id string) (*Dashboard, error)
	ListDashboards(ctx context.Context, organizationID string, offset, limit int) ([]*Dashboard, int, error)
	UpdateDashboard(ctx context.Context, d *Dashboard) error
	DeleteDashboard(ctx context.Context, id string) error
	StoreWidget(ctx context.Context, w *Widget) error
	GetWidget(ctx context.Context, id string) (*Widget, error)
	UpdateWidget(ctx context.Context, w *Widget) error
	DeleteWidget(ctx context.Context, id string) error
	ListWidgets(ctx context.Context, dashboardID string) ([]*Widget, error)
}

type InMemoryRepository struct {
	mu         sync.RWMutex
	dashboards map[string]*Dashboard
	widgets    map[string]*Widget
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		dashboards: make(map[string]*Dashboard),
		widgets:    make(map[string]*Widget),
	}
}

func (r *InMemoryRepository) StoreDashboard(ctx context.Context, d *Dashboard) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now
	r.dashboards[d.ID] = d
	return nil
}

func (r *InMemoryRepository) GetDashboard(ctx context.Context, id string) (*Dashboard, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.dashboards[id]
	if !ok {
		return nil, nil
	}
	var widgets []*Widget
	for _, w := range r.widgets {
		if w.DashboardID == id {
			widgets = append(widgets, w)
		}
	}
	d.Widgets = make([]Widget, len(widgets))
	for i, w := range widgets {
		d.Widgets[i] = *w
	}
	return d, nil
}

func (r *InMemoryRepository) ListDashboards(ctx context.Context, organizationID string, offset, limit int) ([]*Dashboard, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*Dashboard
	for _, d := range r.dashboards {
		if organizationID != "" && d.OrganizationID != organizationID {
			continue
		}
		filtered = append(filtered, d)
	}
	total := len(filtered)
	if offset >= total {
		return []*Dashboard{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return filtered[offset:end], total, nil
}

func (r *InMemoryRepository) UpdateDashboard(ctx context.Context, d *Dashboard) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d.UpdatedAt = time.Now()
	r.dashboards[d.ID] = d
	return nil
}

func (r *InMemoryRepository) DeleteDashboard(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.dashboards, id)
	for wid, w := range r.widgets {
		if w.DashboardID == id {
			delete(r.widgets, wid)
		}
	}
	return nil
}

func (r *InMemoryRepository) StoreWidget(ctx context.Context, w *Widget) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	if w.CreatedAt.IsZero() {
		w.CreatedAt = now
	}
	w.UpdatedAt = now
	r.widgets[w.ID] = w
	return nil
}

func (r *InMemoryRepository) GetWidget(ctx context.Context, id string) (*Widget, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.widgets[id]
	if !ok {
		return nil, nil
	}
	return w, nil
}

func (r *InMemoryRepository) UpdateWidget(ctx context.Context, w *Widget) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	w.UpdatedAt = time.Now()
	r.widgets[w.ID] = w
	return nil
}

func (r *InMemoryRepository) DeleteWidget(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.widgets, id)
	return nil
}

func (r *InMemoryRepository) ListWidgets(ctx context.Context, dashboardID string) ([]*Widget, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Widget
	for _, w := range r.widgets {
		if w.DashboardID == dashboardID {
			result = append(result, w)
		}
	}
	return result, nil
}
