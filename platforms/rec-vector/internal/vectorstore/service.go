package vectorstore

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type VectorStore interface {
	Insert(ctx context.Context, record *VectorRecord) error
	BatchInsert(ctx context.Context, records []*VectorRecord) error
	Delete(ctx context.Context, id, namespace string) error
	Get(ctx context.Context, id, namespace string) (*VectorRecord, error)
	Search(ctx context.Context, query []float64, namespace string, topK int) ([]SearchResult, error)
}

type InMemoryStore struct {
	mu    sync.RWMutex
	store map[string]*VectorRecord
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]*VectorRecord),
	}
}

func recordKey(id, namespace string) string {
	return namespace + ":" + id
}

func (s *InMemoryStore) Insert(ctx context.Context, record *VectorRecord) error {
	if len(record.Vector) == 0 {
		return ErrEmptyVector
	}
	if record.ID == "" {
		record.ID = uuid.New().String()
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[recordKey(record.ID, record.Namespace)] = record
	return nil
}

func (s *InMemoryStore) BatchInsert(ctx context.Context, records []*VectorRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, record := range records {
		if len(record.Vector) == 0 {
			return ErrEmptyVector
		}
		if record.ID == "" {
			record.ID = uuid.New().String()
		}
		if record.CreatedAt.IsZero() {
			record.CreatedAt = time.Now()
		}
		s.store[recordKey(record.ID, record.Namespace)] = record
	}
	return nil
}

func (s *InMemoryStore) Delete(ctx context.Context, id, namespace string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := recordKey(id, namespace)
	if _, ok := s.store[key]; !ok {
		return ErrRecordNotFound
	}
	delete(s.store, key)
	return nil
}

func (s *InMemoryStore) Get(ctx context.Context, id, namespace string) (*VectorRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.store[recordKey(id, namespace)]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return record, nil
}

func (s *InMemoryStore) Search(ctx context.Context, query []float64, namespace string, topK int) ([]SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type scored struct {
		record *VectorRecord
		score  float64
	}
	var candidates []scored

	for _, record := range s.store {
		if record.Namespace == namespace {
			score := cosineSimilarity(query, record.Vector)
			if score > 0 {
				candidates = append(candidates, scored{record: record, score: score})
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	if topK <= 0 {
		topK = len(candidates)
	}
	if topK > len(candidates) {
		topK = len(candidates)
	}

	results := make([]SearchResult, topK)
	for i := 0; i < topK; i++ {
		results[i] = SearchResult{
			ID:        candidates[i].record.ID,
			Score:     math.Round(candidates[i].score*10000) / 10000,
			Metadata:  candidates[i].record.Metadata,
			Namespace: candidates[i].record.Namespace,
			Rank:      i + 1,
		}
	}
	return results, nil
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
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
