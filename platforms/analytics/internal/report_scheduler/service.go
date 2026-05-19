package report_scheduler

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

func (s *Service) ScheduleReport(ctx context.Context, name, description string, query map[string]interface{}, freq ScheduleFrequency, channel DeliveryChannel, recipients []string, webhookURL, format, timeZone, createdBy, orgID string) (*ScheduledReport, error) {
	nextRun := s.calculateNextRun(freq)

	r := &ScheduledReport{
		ID:              uuid.New().String(),
		Name:            name,
		Description:     description,
		Query:           query,
		Frequency:       freq,
		DeliveryChannel: channel,
		Recipients:      recipients,
		WebhookURL:      webhookURL,
		Format:          format,
		TimeZone:        timeZone,
		NextRunAt:       nextRun,
		IsActive:        true,
		CreatedBy:       createdBy,
		OrganizationID:  orgID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.repo.StoreReport(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetReport(ctx context.Context, id string) (*ScheduledReport, error) {
	r, err := s.repo.GetReport(ctx, id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, ErrReportNotFound
	}
	return r, nil
}

func (s *Service) ListReports(ctx context.Context, organizationID string, offset, limit int) ([]*ScheduledReport, int, error) {
	return s.repo.ListReports(ctx, organizationID, offset, limit)
}

func (s *Service) UpdateReport(ctx context.Context, id string, updates map[string]interface{}) (*ScheduledReport, error) {
	r, err := s.repo.GetReport(ctx, id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, ErrReportNotFound
	}
	if v, ok := updates["name"]; ok {
		r.Name = v.(string)
	}
	if v, ok := updates["is_active"]; ok {
		r.IsActive = v.(bool)
	}
	if v, ok := updates["frequency"]; ok {
		r.Frequency = ScheduleFrequency(v.(string))
	}
	r.UpdatedAt = time.Now()
	if err := s.repo.UpdateReport(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteReport(ctx context.Context, id string) error {
	return s.repo.DeleteReport(ctx, id)
}

func (s *Service) GenerateReport(ctx context.Context, reportID string) (*ReportGeneration, error) {
	r, err := s.repo.GetReport(ctx, reportID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, ErrReportNotFound
	}

	gen := &ReportGeneration{
		ID:          uuid.New().String(),
		ReportID:    reportID,
		Status:      "completed",
		GeneratedAt: time.Now(),
	}
	s.repo.StoreGeneration(ctx, gen)

	r.LastRunAt = &gen.GeneratedAt
	r.LastStatus = "completed"
	r.NextRunAt = s.calculateNextRun(r.Frequency)
	s.repo.UpdateReport(ctx, r)

	return gen, nil
}

func (s *Service) DeliverReport(ctx context.Context, genID string, deliveryErr error) error {
	gen, err := s.repo.GetGeneration(ctx, genID)
	if err != nil || gen == nil {
		return ErrReportNotFound
	}
	if deliveryErr != nil {
		gen.Status = "delivery_failed"
		gen.Error = deliveryErr.Error()
	} else {
		gen.Status = "delivered"
		now := time.Now()
		gen.DeliveredAt = &now
	}
	return s.repo.StoreGeneration(ctx, gen)
}

func (s *Service) calculateNextRun(freq ScheduleFrequency) time.Time {
	now := time.Now()
	switch freq {
	case FreqDaily:
		return now.Add(24 * time.Hour).Truncate(24 * time.Hour)
	case FreqWeekly:
		return now.Add(7 * 24 * time.Hour)
	case FreqMonthly:
		return now.AddDate(0, 1, 0)
	default:
		return now.Add(24 * time.Hour)
	}
}
