package vectorstore

import "time"

type VectorRecord struct {
	ID        string                 `json:"id"`
	Vector    []float64              `json:"vector"`
	Metadata  map[string]interface{} `json:"metadata"`
	Namespace string                 `json:"namespace"`
	CreatedAt time.Time              `json:"created_at"`
}

type SearchResult struct {
	ID        string                 `json:"id"`
	Score     float64                `json:"score"`
	Metadata  map[string]interface{} `json:"metadata"`
	Namespace string                 `json:"namespace"`
	Rank      int                    `json:"rank"`
}
