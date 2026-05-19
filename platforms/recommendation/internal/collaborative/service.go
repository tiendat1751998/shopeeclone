package collaborative

import (
	"context"
	"math"
	"sort"
)

type Service interface {
	ItemBasedSimilar(ctx context.Context, itemID string, limit int) ([]ItemSimilarity, error)
	UserBasedRecommend(ctx context.Context, userID string, limit int) ([]ItemSimilarity, error)
	RecordInteraction(ctx context.Context, userID, itemID string, isExplicit bool) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ItemBasedSimilar(ctx context.Context, itemID string, limit int) ([]ItemSimilarity, error) {
	matrix, err := s.repo.GetMatrix(ctx)
	if err != nil {
		return nil, err
	}

	targetRatings, ok := matrix.ItemUserRatings[itemID]
	if !ok {
		return nil, nil
	}

	var similarities []ItemSimilarity
	for otherID, otherRatings := range matrix.ItemUserRatings {
		if otherID == itemID {
			continue
		}
		sim := cosineSimilarity(targetRatings, otherRatings)
		if sim > 0 {
			similarities = append(similarities, ItemSimilarity{
				ItemID:     otherID,
				TargetID:   itemID,
				Similarity: sim,
			})
		}
	}

	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})
	if limit > 0 && limit < len(similarities) {
		similarities = similarities[:limit]
	}
	return similarities, nil
}

func (s *service) UserBasedRecommend(ctx context.Context, userID string, limit int) ([]ItemSimilarity, error) {
	matrix, err := s.repo.GetMatrix(ctx)
	if err != nil {
		return nil, err
	}

	userRatings, ok := matrix.UserItemRatings[userID]
	if !ok {
		return nil, nil
	}

	var userSimilarities []UserSimilarity
	for otherID, otherRatings := range matrix.UserItemRatings {
		if otherID == userID {
			continue
		}
		sim := cosineSimilarity(userRatings, otherRatings)
		if sim > 0 {
			userSimilarities = append(userSimilarities, UserSimilarity{
				UserID:     otherID,
				TargetID:   userID,
				Similarity: sim,
			})
		}
	}

	sort.Slice(userSimilarities, func(i, j int) bool {
		return userSimilarities[i].Similarity > userSimilarities[j].Similarity
	})
	if len(userSimilarities) > 20 {
		userSimilarities = userSimilarities[:20]
	}

	candidateScores := make(map[string]float64)
	for _, us := range userSimilarities {
		otherRatings := matrix.UserItemRatings[us.UserID]
		for itemID, rating := range otherRatings {
			if _, already := userRatings[itemID]; already {
				continue
			}
			candidateScores[itemID] += us.Similarity * rating
		}
	}

	var recommendations []ItemSimilarity
	for itemID, score := range candidateScores {
		recommendations = append(recommendations, ItemSimilarity{
			ItemID:     itemID,
			TargetID:   userID,
			Similarity: score,
		})
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Similarity > recommendations[j].Similarity
	})
	if limit > 0 && limit < len(recommendations) {
		recommendations = recommendations[:limit]
	}
	return recommendations, nil
}

func (s *service) RecordInteraction(ctx context.Context, userID, itemID string, isExplicit bool) error {
	weight := 0.3
	if isExplicit {
		weight = 1.0
	}
	return s.repo.StoreRating(ctx, userID, itemID, weight)
}

func cosineSimilarity(a, b map[string]float64) float64 {
	var dotProduct, normA, normB float64
	for k, v := range a {
		dotProduct += v * b[k]
		normA += v * v
	}
	for _, v := range b {
		normB += v * v
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
