package embeddings

import "time"

type Embedding struct {
	EntityID    string    `json:"entity_id"`
	EntityType  string    `json:"entity_type"`
	Vector      []float64 `json:"vector"`
	ModelName   string    `json:"model_name"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
}

type SimilarityResult struct {
	EntityID   string  `json:"entity_id"`
	EntityType string  `json:"entity_type"`
	Score      float64 `json:"score"`
}
