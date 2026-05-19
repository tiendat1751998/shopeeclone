package userembedding

import "time"

type UserEmbedding struct {
	UserID       string    `json:"user_id"`
	Embedding    []float64 `json:"embedding"`
	UpdatedAt    time.Time `json:"updated_at"`
	ModelVersion string    `json:"model_version"`
}
