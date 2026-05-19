package reranker

import "time"

type ReRankedResult struct {
	ProductID string  `json:"product_id"`
	OriginalScore float64 `json:"original_score"`
	AdjustedScore  float64 `json:"adjusted_score"`
	Category       string  `json:"category"`
	CreatedAt      time.Time `json:"created_at"`
	ExposureCount  int     `json:"exposure_count"`
}

type BoostFactor struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
