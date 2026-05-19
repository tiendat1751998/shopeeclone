package reporting

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("reporting: not found")

type Service interface {
	TrackSend(ctx context.Context, campaignID string) error
	TrackOpen(ctx context.Context, campaignID string) error
	TrackClick(ctx context.Context, campaignID string) error
	TrackConversion(ctx context.Context, req *TrackEventRequest) error
	TrackBounce(ctx context.Context, campaignID string) error
	TrackUnsubscribe(ctx context.Context, campaignID string) error
	GetCampaignReport(ctx context.Context, campaignID string) (*CampaignReport, error)
	GetAggregatedReport(ctx context.Context) (*AggregatedReport, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) getOrCreateReport(ctx context.Context, campaignID string) (*CampaignReport, error) {
	report, err := s.repo.GetReport(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	if report == nil {
		report = &CampaignReport{CampaignID: campaignID}
	}
	return report, nil
}

func (s *service) TrackSend(ctx context.Context, campaignID string) error {
	report, err := s.getOrCreateReport(ctx, campaignID)
	if err != nil {
		return err
	}
	report.SentCount++
	return s.repo.UpsertReport(ctx, report)
}

func (s *service) TrackOpen(ctx context.Context, campaignID string) error {
	report, err := s.getOrCreateReport(ctx, campaignID)
	if err != nil {
		return err
	}
	report.OpenedCount++
	return s.repo.UpsertReport(ctx, report)
}

func (s *service) TrackClick(ctx context.Context, campaignID string) error {
	report, err := s.getOrCreateReport(ctx, campaignID)
	if err != nil {
		return err
	}
	report.ClickedCount++
	return s.repo.UpsertReport(ctx, report)
}

func (s *service) TrackConversion(ctx context.Context, req *TrackEventRequest) error {
	report, err := s.getOrCreateReport(ctx, req.CampaignID)
	if err != nil {
		return err
	}
	report.ConvertedCount++
	report.RevenueAttributed += req.Revenue
	return s.repo.UpsertReport(ctx, report)
}

func (s *service) TrackBounce(ctx context.Context, campaignID string) error {
	report, err := s.getOrCreateReport(ctx, campaignID)
	if err != nil {
		return err
	}
	report.BouncedCount++
	return s.repo.UpsertReport(ctx, report)
}

func (s *service) TrackUnsubscribe(ctx context.Context, campaignID string) error {
	report, err := s.getOrCreateReport(ctx, campaignID)
	if err != nil {
		return err
	}
	report.UnsubscribedCount++
	return s.repo.UpsertReport(ctx, report)
}

func (s *service) GetCampaignReport(ctx context.Context, campaignID string) (*CampaignReport, error) {
	report, err := s.repo.GetReport(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, ErrNotFound
	}
	return report, nil
}

func (s *service) GetAggregatedReport(ctx context.Context) (*AggregatedReport, error) {
	reports, err := s.repo.ListReports(ctx)
	if err != nil {
		return nil, err
	}

	agg := &AggregatedReport{}
	for _, r := range reports {
		agg.TotalCampaigns++
		agg.TotalSent += r.SentCount
		agg.TotalDelivered += r.DeliveredCount
		agg.TotalOpened += r.OpenedCount
		agg.TotalClicked += r.ClickedCount
		agg.TotalConverted += r.ConvertedCount
		agg.TotalBounced += r.BouncedCount
		agg.TotalUnsubscribed += r.UnsubscribedCount
		agg.TotalRevenue += r.RevenueAttributed
	}

	if agg.TotalDelivered > 0 {
		agg.OverallOpenRate = float64(agg.TotalOpened) / float64(agg.TotalDelivered) * 100
		agg.OverallClickRate = float64(agg.TotalClicked) / float64(agg.TotalDelivered) * 100
		agg.OverallConversionRate = float64(agg.TotalConverted) / float64(agg.TotalDelivered) * 100
	}

	return agg, nil
}
