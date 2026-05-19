package similarity

type SimilarityRequest struct {
	QueryEmbedding []float64              `json:"query_embedding"`
	Namespace      string                 `json:"namespace"`
	TopK           int                    `json:"top_k"`
	MinScore       float64                `json:"min_score"`
	Filter         map[string]interface{} `json:"filter"`
	Keyword        string                 `json:"keyword"`
}

type SimilarityResult struct {
	ID        string                 `json:"id"`
	Score     float64                `json:"score"`
	Metadata  map[string]interface{} `json:"metadata"`
	Rank      int                    `json:"rank"`
}
