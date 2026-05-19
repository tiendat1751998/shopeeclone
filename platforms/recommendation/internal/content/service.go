package content

import (
	"context"
	"math"
	"sort"
	"strings"
)

type Service interface {
	SimilarByContent(ctx context.Context, productID string, limit int) ([]string, error)
	ScoreProduct(ctx context.Context, source, target *ProductFeatures) (float64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SimilarByContent(ctx context.Context, productID string, limit int) ([]string, error) {
	source, err := s.repo.GetProductFeatures(ctx, productID)
	if err != nil || source == nil {
		return nil, err
	}

	all, err := s.repo.GetAllProductFeatures(ctx)
	if err != nil {
		return nil, err
	}

	type scored struct {
		productID string
		score     float64
	}
	var results []scored

	for _, target := range all {
		if target.ProductID == productID {
			continue
		}
		score, err := s.ScoreProduct(ctx, source, &target)
		if err != nil {
			continue
		}
		if score > 0 {
			results = append(results, scored{productID: target.ProductID, score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}

	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r.productID
	}
	return ids, nil
}

func (s *service) ScoreProduct(ctx context.Context, source, target *ProductFeatures) (float64, error) {
	var score float64

	categoryScore := categoryOverlapScore(source.Category, source.ParentCategory, target.Category, target.ParentCategory)
	score += categoryScore * 0.4

	tagScore := jaccardSimilarity(source.Tags, target.Tags)
	score += tagScore * 0.35

	priceScore := priceProximityScore(source.Price, target.Price)
	score += priceScore * 0.25

	return score, nil
}

func categoryOverlapScore(catA, parentA, catB, parentB string) float64 {
	catA = strings.ToLower(strings.TrimSpace(catA))
	parentA = strings.ToLower(strings.TrimSpace(parentA))
	catB = strings.ToLower(strings.TrimSpace(catB))
	parentB = strings.ToLower(strings.TrimSpace(parentB))

	if catA != "" && catA == catB {
		return 1.0
	}
	if parentA != "" && parentA == parentB {
		return 0.5
	}
	if catA != "" && parentB != "" && catA == parentB {
		return 0.5
	}
	if parentA != "" && catB != "" && parentA == catB {
		return 0.5
	}
	return 0
}

func jaccardSimilarity(a, b []string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 0
	}
	setA := make(map[string]bool)
	for _, t := range a {
		setA[strings.ToLower(t)] = true
	}
	setB := make(map[string]bool)
	for _, t := range b {
		setB[strings.ToLower(t)] = true
	}
	intersection := 0
	for t := range setA {
		if setB[t] {
			intersection++
		}
	}
	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

func priceProximityScore(priceA, priceB float64) float64 {
	if priceA <= 0 || priceB <= 0 {
		return 0
	}
	ratio := priceA / priceB
	if ratio > 1 {
		ratio = 1 / ratio
	}
	score := 1 - math.Abs(1-ratio)
	if score < 0 {
		score = 0
	}
	return score
}
