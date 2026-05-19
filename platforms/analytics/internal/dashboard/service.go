package dashboard

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateDashboard(ctx context.Context, title, description, orgID, createdBy string, isPublic bool, tags []string) (*Dashboard, error) {
	d := &Dashboard{
		ID:             uuid.New().String(),
		Title:          title,
		Description:    description,
		OrganizationID: orgID,
		CreatedBy:      createdBy,
		IsPublic:       isPublic,
		Tags:           tags,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.StoreDashboard(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetDashboard(ctx context.Context, id string) (*Dashboard, error) {
	d, err := s.repo.GetDashboard(ctx, id)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, ErrDashboardNotFound
	}
	return d, nil
}

func (s *Service) ListDashboards(ctx context.Context, organizationID string, offset, limit int) ([]*Dashboard, int, error) {
	return s.repo.ListDashboards(ctx, organizationID, offset, limit)
}

func (s *Service) UpdateDashboard(ctx context.Context, id, title, description string, isPublic bool, tags []string) (*Dashboard, error) {
	d, err := s.repo.GetDashboard(ctx, id)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, ErrDashboardNotFound
	}
	d.Title = title
	d.Description = description
	d.IsPublic = isPublic
	d.Tags = tags
	d.UpdatedAt = time.Now()
	if err := s.repo.UpdateDashboard(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) DeleteDashboard(ctx context.Context, id string) error {
	return s.repo.DeleteDashboard(ctx, id)
}

func (s *Service) AddWidget(ctx context.Context, dashboardID, title string, chartType string, width, height, posX, posY int, dataSource DataSource, config map[string]interface{}) (*Widget, error) {
	d, err := s.repo.GetDashboard(ctx, dashboardID)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, ErrDashboardNotFound
	}

	w := &Widget{
		ID:          uuid.New().String(),
		DashboardID: dashboardID,
		Title:       title,
		Type:        ChartType(chartType),
		Width:       width,
		Height:      height,
		PositionX:   posX,
		PositionY:   posY,
		DataSource:  dataSource,
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.StoreWidget(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) UpdateWidget(ctx context.Context, widgetID, title string, chartType string, width, height, posX, posY int, dataSource DataSource, config map[string]interface{}) (*Widget, error) {
	w, err := s.repo.GetWidget(ctx, widgetID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrWidgetNotFound
	}
	w.Title = title
	w.Type = ChartType(chartType)
	w.Width = width
	w.Height = height
	w.PositionX = posX
	w.PositionY = posY
	w.DataSource = dataSource
	w.Config = config
	w.UpdatedAt = time.Now()
	if err := s.repo.UpdateWidget(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) RemoveWidget(ctx context.Context, widgetID string) error {
	return s.repo.DeleteWidget(ctx, widgetID)
}

func (s *Service) RefreshWidget(ctx context.Context, widgetID string) (*Widget, error) {
	w, err := s.repo.GetWidget(ctx, widgetID)
	if err != nil || w == nil {
		return nil, ErrWidgetNotFound
	}
	return w, nil
}
