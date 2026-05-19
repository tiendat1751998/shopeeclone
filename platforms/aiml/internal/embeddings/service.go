package embeddings

import (
	"context"
	"hash/fnv"
	"math"
	"sync"
	"time"
)

const DefaultDimension = 128

type EmbeddingGenerator struct{}

func NewEmbeddingGenerator() *EmbeddingGenerator {
	return &EmbeddingGenerator{}
}

func (g *EmbeddingGenerator) Generate(entityID, entityType, modelName, version string) *Embedding {
	vector := make([]float64, DefaultDimension)
	h := fnv.New64a()
	for i := 0; i < DefaultDimension; i++ {
		h.Reset()
		h.Write([]byte(entityID + ":" + entityType + ":" + modelName + ":" + version + ":" + string(rune(i))))
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
	return &Embedding{
		EntityID:   entityID,
		EntityType: entityType,
		Vector:     vector,
		ModelName:  modelName,
		Version:    version,
		CreatedAt:  time.Now(),
	}
}

type VectorStore interface {
	Store(ctx context.Context, emb *Embedding) error
	Get(ctx context.Context, entityID, entityType string) (*Embedding, error)
	FindSimilar(ctx context.Context, vector []float64, entityType string, topK int) ([]SimilarityResult, error)
}

type InMemoryVectorStore struct {
	mu         sync.RWMutex
	embeddings map[string]*Embedding
}

func NewInMemoryVectorStore() *InMemoryVectorStore {
	return &InMemoryVectorStore{
		embeddings: make(map[string]*Embedding),
	}
}

func embeddingKey(entityID, entityType string) string {
	return entityType + ":" + entityID
}

func (s *InMemoryVectorStore) Store(ctx context.Context, emb *Embedding) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if emb.CreatedAt.IsZero() {
		emb.CreatedAt = time.Now()
	}
	s.embeddings[embeddingKey(emb.EntityID, emb.EntityType)] = emb
	return nil
}

func (s *InMemoryVectorStore) Get(ctx context.Context, entityID, entityType string) (*Embedding, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	emb, ok := s.embeddings[embeddingKey(entityID, entityType)]
	if !ok {
		return nil, ErrEmbeddingNotFound
	}
	return emb, nil
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (s *InMemoryVectorStore) FindSimilar(ctx context.Context, vector []float64, entityType string, topK int) ([]SimilarityResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type scored struct {
		entityID   string
		entityType string
		score      float64
	}
	var scoredList []scored

	for key, emb := range s.embeddings {
		if entityType != "" && emb.EntityType != entityType {
			continue
		}
		_ = key
		score := cosineSimilarity(vector, emb.Vector)
		if score > 0 {
			scoredList = append(scoredList, scored{
				entityID:   emb.EntityID,
				entityType: emb.EntityType,
				score:      score,
			})
		}
	}

	for i := 0; i < len(scoredList); i++ {
		for j := i + 1; j < len(scoredList); j++ {
			if scoredList[j].score > scoredList[i].score {
				scoredList[i], scoredList[j] = scoredList[j], scoredList[i]
			}
		}
	}

	if topK <= 0 || topK > len(scoredList) {
		topK = len(scoredList)
	}

	results := make([]SimilarityResult, topK)
	for i := 0; i < topK; i++ {
		results[i] = SimilarityResult{
			EntityID:   scoredList[i].entityID,
			EntityType: scoredList[i].entityType,
			Score:      math.Round(scoredList[i].score*10000) / 10000,
		}
	}
	return results, nil
}

type Service struct {
	generator *EmbeddingGenerator
	store     VectorStore
}

func NewService(generator *EmbeddingGenerator, store VectorStore) *Service {
	return &Service{generator: generator, store: store}
}

func (s *Service) GenerateEmbedding(ctx context.Context, entityID, entityType, modelName, version string) (*Embedding, error) {
	emb := s.generator.Generate(entityID, entityType, modelName, version)
	if err := s.store.Store(ctx, emb); err != nil {
		return nil, err
	}
	return emb, nil
}

func (s *Service) GetEmbedding(ctx context.Context, entityID, entityType string) (*Embedding, error) {
	return s.store.Get(ctx, entityID, entityType)
}

func (s *Service) FindSimilar(ctx context.Context, entityID, entityType string, topK int) ([]SimilarityResult, error) {
	emb, err := s.store.Get(ctx, entityID, entityType)
	if err != nil {
		return nil, err
	}
	return s.store.FindSimilar(ctx, emb.Vector, entityType, topK)
}

func (s *Service) FindSimilarByVector(ctx context.Context, vector []float64, entityType string, topK int) ([]SimilarityResult, error) {
	return s.store.FindSimilar(ctx, vector, entityType, topK)
}

func (s *Service) StoreEmbedding(ctx context.Context, emb *Embedding) error {
	return s.store.Store(ctx, emb)
}

func (s *Service) BatchGenerate(ctx context.Context, entityIDs []string, entityType, modelName, version string) ([]*Embedding, error) {
	embeddings := make([]*Embedding, 0, len(entityIDs))
	for _, id := range entityIDs {
		emb := s.generator.Generate(id, entityType, modelName, version)
		if err := s.store.Store(ctx, emb); err != nil {
			return embeddings, err
		}
		embeddings = append(embeddings, emb)
	}
	return embeddings, nil
}
