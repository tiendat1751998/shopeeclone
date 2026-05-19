package monitoring

import (
	"context"
	"time"
)

type Service interface {
	ReportMetrics(ctx context.Context, metric *IndexMetric) error
	GetIndexMetrics(ctx context.Context) ([]*IndexMetric, error)
	GetIndexMetric(ctx context.Context, indexName string) (*IndexMetric, error)
	GetClusterHealth(ctx context.Context) (*ClusterHealth, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ReportMetrics(ctx context.Context, metric *IndexMetric) error {
	if metric.Health == "" {
		metric.Health = HealthGreen
	}
	metric.LastUpdated = time.Now()
	return s.repo.UpsertMetric(ctx, metric)
}

func (s *service) GetIndexMetrics(ctx context.Context) ([]*IndexMetric, error) {
	return s.repo.ListMetrics(ctx)
}

func (s *service) GetIndexMetric(ctx context.Context, indexName string) (*IndexMetric, error) {
	return s.repo.GetMetric(ctx, indexName)
}

func (s *service) GetClusterHealth(ctx context.Context) (*ClusterHealth, error) {
	metrics, err := s.repo.ListMetrics(ctx)
	if err != nil {
		return nil, err
	}

	health := &ClusterHealth{}
	for _, m := range metrics {
		health.ActiveShards += m.ShardCount
		if m.Health == HealthRed {
			health.UnassignedShards++
		} else if m.Health == HealthYellow {
			health.RelocatingShards++
		}
	}
	health.NodesCount = len(metrics)
	if health.NodesCount == 0 {
		health.NodesCount = 1
	}

	return health, nil
}
