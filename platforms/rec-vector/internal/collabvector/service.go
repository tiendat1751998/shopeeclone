package collabvector

import (
	"context"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordInteraction(ctx context.Context, userID, itemID, interactionType string) error {
	weight := implicitWeight(interactionType)
	interaction := &Interaction{
		UserID:          userID,
		ItemID:          itemID,
		InteractionType: interactionType,
		Weight:          weight,
		Timestamp:       time.Now(),
	}
	return s.repo.StoreInteraction(ctx, interaction)
}

func (s *Service) TrainFactorization(ctx context.Context, numFactors, numIterations int, learningRate float64) error {
	matrix, err := s.repo.GetMatrix(ctx)
	if err != nil {
		return err
	}
	if len(matrix.Users) == 0 || len(matrix.Items) == 0 {
		return ErrNotEnoughInteractions
	}

	numUsers := len(matrix.Users)
	numItems := len(matrix.Items)

	userFactors := make([][]float64, numUsers)
	itemFactors := make([][]float64, numItems)
	for i := 0; i < numUsers; i++ {
		userFactors[i] = make([]float64, numFactors)
		for j := 0; j < numFactors; j++ {
			userFactors[i][j] = (rand.Float64() - 0.5) * 0.1
		}
	}
	for i := 0; i < numItems; i++ {
		itemFactors[i] = make([]float64, numFactors)
		for j := 0; j < numFactors; j++ {
			itemFactors[i][j] = (rand.Float64() - 0.5) * 0.1
		}
	}

	for iter := 0; iter < numIterations; iter++ {
		for u := 0; u < numUsers; u++ {
			for i := 0; i < numItems; i++ {
				rating, ok := matrix.UserItemRatings[matrix.Users[u]][matrix.Items[i]]
				if !ok {
					continue
				}
				pred := dotProduct(userFactors[u], itemFactors[i])
				err := rating - pred
				for f := 0; f < numFactors; f++ {
					userFactors[u][f] += learningRate * (err * itemFactors[i][f])
					itemFactors[i][f] += learningRate * (err * userFactors[u][f])
				}
			}
		}
	}

	return s.repo.StoreFactors(ctx, userFactors, itemFactors)
}

func (s *Service) GenerateLatentFactors(ctx context.Context, numFactors int) (*LatentFactors, error) {
	if err := s.TrainFactorization(ctx, numFactors, 20, 0.01); err != nil {
		return nil, err
	}
	return s.repo.GetFactors(ctx)
}

func (s *Service) RecommendByFactorization(ctx context.Context, userID string, topK int) ([]FactorRecommendation, error) {
	matrix, err := s.repo.GetMatrix(ctx)
	if err != nil {
		return nil, err
	}
	userIdx, ok := matrix.UserIndex[userID]
	if !ok {
		return nil, ErrUserNotFound
	}

	factors, err := s.repo.GetFactors(ctx)
	if err != nil {
		return nil, err
	}

	userRatings := matrix.UserItemRatings[userID]

	scores := make([]FactorRecommendation, 0, len(matrix.Items))
	for i, itemID := range matrix.Items {
		if _, rated := userRatings[itemID]; rated {
			continue
		}
		score := dotProduct(factors.UserFactors[userIdx], factors.ItemFactors[i])
		scores = append(scores, FactorRecommendation{
			ItemID: itemID,
			Score:  math.Round(score*10000) / 10000,
		})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	if topK <= 0 || topK > len(scores) {
		topK = len(scores)
	}
	return scores[:topK], nil
}

func implicitWeight(interactionType string) float64 {
	switch interactionType {
	case "purchase":
		return 1.0
	case "click":
		return 0.5
	case "view":
		return 0.1
	default:
		return 0.1
	}
}

func dotProduct(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}
