package userembedding

import (
	"context"
	"hash/fnv"
	"math"
	"time"
)

const DefaultDimension = 128

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func generateHashVector(seed string, dim int) []float64 {
	vector := make([]float64, dim)
	h := fnv.New64a()
	for i := 0; i < dim; i++ {
		h.Reset()
		h.Write([]byte(seed + ":" + string(rune(i))))
		vector[i] = (float64(h.Sum64() % 200000) / 100000.0) - 1.0
	}
	norm := 0.0
	for _, v := range vector {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range vector {
			vector[i] /= norm
		}
	}
	return vector
}

func (s *Service) GenerateUserEmbedding(ctx context.Context, userID string, modelVersion string) (*UserEmbedding, error) {
	emb := &UserEmbedding{
		UserID:       userID,
		Embedding:    generateHashVector("user:"+userID+":"+modelVersion, DefaultDimension),
		UpdatedAt:    time.Now(),
		ModelVersion: modelVersion,
	}
	if err := s.repo.Store(ctx, emb); err != nil {
		return nil, err
	}
	return emb, nil
}

func (s *Service) UpdateEmbedding(ctx context.Context, userID string, embedding []float64, modelVersion string) (*UserEmbedding, error) {
	emb := &UserEmbedding{
		UserID:       userID,
		Embedding:    embedding,
		UpdatedAt:    time.Now(),
		ModelVersion: modelVersion,
	}
	if err := s.repo.Store(ctx, emb); err != nil {
		return nil, err
	}
	return emb, nil
}

func (s *Service) GetEmbedding(ctx context.Context, userID string) (*UserEmbedding, error) {
	return s.repo.Get(ctx, userID)
}

func (s *Service) BatchGetEmbeddings(ctx context.Context, userIDs []string) ([]*UserEmbedding, error) {
	results := make([]*UserEmbedding, 0, len(userIDs))
	for _, uid := range userIDs {
		emb, err := s.repo.Get(ctx, uid)
		if err != nil {
			continue
		}
		results = append(results, emb)
	}
	return results, nil
}
