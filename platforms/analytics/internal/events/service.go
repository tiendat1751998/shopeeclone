package events

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo     Repository
	dedupMap sync.Map
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) IngestEvent(ctx context.Context, event *AnalyticsEvent) error {
	if event.EventID == "" {
		event.EventID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	if _, loaded := s.dedupMap.LoadOrStore(event.EventID, true); loaded {
		return nil
	}

	return s.repo.StoreEvent(ctx, event)
}

func (s *Service) BatchIngest(ctx context.Context, events []AnalyticsEvent) (int, error) {
	ingested := 0
	for i := range events {
		if err := s.IngestEvent(ctx, &events[i]); err != nil {
			return ingested, err
		}
		ingested++
	}
	return ingested, nil
}

func (s *Service) ReplayEvents(ctx context.Context, startTime, endTime time.Time, handler func(*AnalyticsEvent) error) error {
	events_list, _, err := s.repo.ListEvents(ctx, "", startTime, endTime, 0, 10000)
	if err != nil {
		return err
	}
	for _, event := range events_list {
		if err := handler(event); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetEvent(ctx context.Context, eventID string) (*AnalyticsEvent, error) {
	return s.repo.GetEvent(ctx, eventID)
}

func (s *Service) ListEvents(ctx context.Context, eventType EventType, startTime, endTime time.Time, offset, limit int) ([]*AnalyticsEvent, int64, error) {
	return s.repo.ListEvents(ctx, eventType, startTime, endTime, offset, limit)
}

func (s *Service) GetEventCount(ctx context.Context, eventType EventType, startTime, endTime time.Time) (int64, error) {
	return s.repo.GetEventCount(ctx, eventType, startTime, endTime)
}

func (s *Service) GetRevenue(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	return s.repo.GetRevenue(ctx, startTime, endTime)
}

func (s *Service) GetUniqueUsers(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	return s.repo.GetUniqueUsers(ctx, startTime, endTime)
}

func (s *Service) GetActiveUsers(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	return s.repo.GetActiveUsers(ctx, startTime, endTime)
}

func (s *Service) GetOrders(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	return s.repo.GetOrders(ctx, startTime, endTime)
}
