package itemembedding

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

func (s *Service) GenerateItemEmbedding(ctx context.Context, itemID, category string, tags []string, modelVersion string) (*ItemEmbedding, error) {
	emb := &ItemEmbedding{
		ItemID:       itemID,
		Embedding:    generateHashVector("item:"+itemID+":"+category+":"+modelVersion, DefaultDimension),
		Category:     category,
		Tags:         tags,
		UpdatedAt:    time.Now(),
		ModelVersion: modelVersion,
	}
	if err := s.repo.Store(ctx, emb); err != nil {
		return nil, err
	}
	return emb, nil
}

func (s *Service) UpdateEmbedding(ctx context.Context, itemID, category string, tags []string, embedding []float64, modelVersion string) (*ItemEmbedding, error) {
	emb := &ItemEmbedding{
		ItemID:       itemID,
		Embedding:    embedding,
		Category:     category,
		Tags:         tags,
		UpdatedAt:    time.Now(),
		ModelVersion: modelVersion,
	}
	if err := s.repo.Store(ctx, emb); err != nil {
		return nil, err
	}
	return emb, nil
}

func (s *Service) GetEmbedding(ctx context.Context, itemID string) (*ItemEmbedding, error) {
	return s.repo.Get(ctx, itemID)
}

func (s *Service) BatchGetEmbeddings(ctx context.Context, itemIDs []string) ([]*ItemEmbedding, error) {
	results := make([]*ItemEmbedding, 0, len(itemIDs))
	for _, id := range itemIDs {
		emb, err := s.repo.Get(ctx, id)
		if err != nil {
			continue
		}
		results = append(results, emb)
	}
	return results, nil
}
