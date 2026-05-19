package personalization

import (
	"context"
	"math"
	"time"
)

type InteractionEvent struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Category  string    `json:"category"`
	Brand     string    `json:"brand"`
	Price     float64   `json:"price"`
	Tags      []string  `json:"tags"`
	Weight    float64   `json:"weight"`
	Timestamp time.Time `json:"timestamp"`
}

type Service interface {
	GetProfile(ctx context.Context, userID string) (*UserProfile, error)
	BuildProfile(ctx context.Context, events []InteractionEvent) (*UserProfile, error)
	ScoreItem(ctx context.Context, profile *UserProfile, category string, brand string, price float64, tags []string) float64
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}

func (s *service) BuildProfile(ctx context.Context, events []InteractionEvent) (*UserProfile, error) {
	if len(events) == 0 {
		return nil, nil
	}

	profile := &UserProfile{
		UserID:          events[0].UserID,
		CategoryWeights: make(map[string]float64),
		PreferredBrands: make(map[string]float64),
		InterestVector:  make(map[string]float64),
		PriceRangeMin:   math.MaxFloat64,
		PriceRangeMax:   0,
	}

	now := time.Now()
	totalWeight := 0.0
	totalPrice := 0.0

	for _, evt := range events {
		hoursAge := now.Sub(evt.Timestamp).Hours()
		decay := math.Exp(-hoursAge / 720)

		weight := evt.Weight * decay
		totalWeight += weight

		if evt.Category != "" {
			profile.CategoryWeights[evt.Category] += weight
		}
		if evt.Brand != "" {
			profile.PreferredBrands[evt.Brand] += weight
		}
		if evt.Price > 0 {
			if evt.Price < profile.PriceRangeMin {
				profile.PriceRangeMin = evt.Price
			}
			if evt.Price > profile.PriceRangeMax {
				profile.PriceRangeMax = evt.Price
			}
			totalPrice += evt.Price
		}
		for _, tag := range evt.Tags {
			profile.InterestVector[tag] += weight
		}
	}

	if totalWeight > 0 {
		for cat := range profile.CategoryWeights {
			profile.CategoryWeights[cat] = math.Round(profile.CategoryWeights[cat]/totalWeight*1000) / 1000
		}
		for brand := range profile.PreferredBrands {
			profile.PreferredBrands[brand] = math.Round(profile.PreferredBrands[brand]/totalWeight*1000) / 1000
		}
		for tag := range profile.InterestVector {
			profile.InterestVector[tag] = math.Round(profile.InterestVector[tag]/totalWeight*1000) / 1000
		}
	}

	if len(events) > 0 {
		profile.PreferredPriceMid = math.Round(totalPrice/float64(len(events))*100) / 100
	}

	if profile.PriceRangeMin == math.MaxFloat64 {
		profile.PriceRangeMin = 0
	}

	profile.TotalInteractions = len(events)

	if err := s.repo.SaveProfile(ctx, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *service) ScoreItem(ctx context.Context, profile *UserProfile, category string, brand string, price float64, tags []string) float64 {
	if profile == nil {
		return 0
	}

	var score float64

	if category != "" {
		if catWeight, ok := profile.CategoryWeights[category]; ok {
			score += catWeight * 0.4
		}
	}

	if brand != "" {
		if brandWeight, ok := profile.PreferredBrands[brand]; ok {
			score += brandWeight * 0.25
		}
	}

	if price > 0 && profile.PreferredPriceMid > 0 {
		ratio := price / profile.PreferredPriceMid
		if ratio > 1 {
			ratio = 1 / ratio
		}
		priceScore := 1 - math.Abs(1-ratio)
		if priceScore > 0 {
			score += priceScore * 0.2
		}
	}

	if len(tags) > 0 {
		tagMatchCount := 0
		for _, tag := range tags {
			if _, ok := profile.InterestVector[tag]; ok {
				tagMatchCount++
			}
		}
		if tagMatchCount > 0 {
			tagScore := float64(tagMatchCount) / float64(len(tags))
			score += tagScore * 0.15
		}
	}

	return math.Round(score*1000) / 1000
}


