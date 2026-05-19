package itemembedding

import "time"

type ItemEmbedding struct {
	ItemID       string    `json:"item_id"`
	Embedding    []float64 `json:"embedding"`
	Category     string    `json:"category"`
	Tags         []string  `json:"tags"`
	UpdatedAt    time.Time `json:"updated_at"`
	ModelVersion string    `json:"model_version"`
}
