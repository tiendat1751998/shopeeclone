package trending

import (
	"context"
	"math"
	"sort"
	"time"
)

type Service interface {
	GetTrending(ctx context.Context, limit int) ([]TrendingScore, error)
	RecordInteraction(ctx context.Context, productID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) getWeight(age time.Duration) float64 {
	if age <= 1*time.Hour {
		return 1.0
	}
	if age <= 24*time.Hour {
		return 0.5
	}
	if age <= 7*24*time.Hour {
		return 0.1
	}
	return 0
}

func (s *service) GetTrending(ctx context.Context, limit int) ([]TrendingScore, error) {
	now := time.Now()
	windowStart := now.Add(-7 * 24 * time.Hour)

	interactions, err := s.repo.GetWindowInteractions(ctx, windowStart)
	if err != nil {
		return nil, err
	}

	productScores := make(map[string]float64)
	productCounts := make(map[string]int)
	maxScore := 0.0

	for _, inter := range interactions {
		age := now.Sub(inter.Timestamp)
		weight := s.getWeight(age)
		score := inter.Weight * weight
		productScores[inter.ProductID] += score
		productCounts[inter.ProductID]++
		if productScores[inter.ProductID] > maxScore {
			maxScore = productScores[inter.ProductID]
		}
	}

	var results []TrendingScore
	for pid, score := range productScores {
		normalized := 0.0
		if maxScore > 0 {
			normalized = score / maxScore
		}
		velocity := 0.0
		if count, ok := productCounts[pid]; ok && count > 0 {
			velocity = float64(count) / 7.0
		}
		results = append(results, TrendingScore{
			ProductID: pid,
			Score:     math.Round(normalized*1000) / 1000,
			Velocity:  math.Round(velocity*100) / 100,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}
	return results, nil
}

func (s *service) RecordInteraction(ctx context.Context, productID string) error {
	return s.repo.RecordInteraction(ctx, productID, 1.0)
}
