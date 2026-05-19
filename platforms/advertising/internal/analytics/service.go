package analytics

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	RecordImpression(ctx context.Context, imp *Impression) error
	RecordClick(ctx context.Context, click *Click) error
	RecordConversion(ctx context.Context, conv *Conversion) error
	GetReport(ctx context.Context, filter *ReportFilter) (*AnalyticsReport, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RecordImpression(ctx context.Context, imp *Impression) error {
	imp.ID = uuid.New().String()
	if imp.Timestamp.IsZero() {
		imp.Timestamp = time.Now()
	}
	return s.repo.StoreImpression(ctx, imp)
}

func (s *service) RecordClick(ctx context.Context, click *Click) error {
	click.ID = uuid.New().String()
	if click.Timestamp.IsZero() {
		click.Timestamp = time.Now()
	}
	return s.repo.StoreClick(ctx, click)
}

func (s *service) RecordConversion(ctx context.Context, conv *Conversion) error {
	conv.ID = uuid.New().String()
	if conv.Timestamp.IsZero() {
		conv.Timestamp = time.Now()
	}
	return s.repo.StoreConversion(ctx, conv)
}

func (s *service) GetReport(ctx context.Context, filter *ReportFilter) (*AnalyticsReport, error) {
	impressions, err := s.repo.GetImpressions(ctx, filter)
	if err != nil {
		return nil, err
	}
	clicks, err := s.repo.GetClicks(ctx, filter)
	if err != nil {
		return nil, err
	}
	conversions, err := s.repo.GetConversions(ctx, filter)
	if err != nil {
		return nil, err
	}

	report := &AnalyticsReport{
		CampaignID:  filter.CampaignID,
		Impressions: int64(len(impressions)),
		Clicks:      int64(len(clicks)),
		Conversions: int64(len(conversions)),
	}

	for _, imp := range impressions {
		report.Spend += imp.Cost
	}
	for _, cl := range clicks {
		report.Spend += cl.Cost
	}
	for _, conv := range conversions {
		report.Revenue += conv.Revenue
	}

	if report.Impressions > 0 {
		report.CTR = float64(report.Clicks) / float64(report.Impressions) * 100
		report.CPM = (report.Spend / float64(report.Impressions)) * 1000
	}
	if report.Clicks > 0 {
		report.CVR = float64(report.Conversions) / float64(report.Clicks) * 100
		report.CPC = report.Spend / float64(report.Clicks)
	}
	if report.Spend > 0 {
		report.ROAS = report.Revenue / report.Spend
	}

	report.CTR = math.Round(report.CTR*100) / 100
	report.CVR = math.Round(report.CVR*100) / 100
	report.CPC = math.Round(report.CPC*100) / 100
	report.CPM = math.Round(report.CPM*100) / 100
	report.ROAS = math.Round(report.ROAS*100) / 100

	return report, nil
}
