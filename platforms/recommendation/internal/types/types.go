package types

import "time"

type RecommendationType string

const (
	RecTypeRelated     RecommendationType = "related"
	RecTypeTrending    RecommendationType = "trending"
	RecTypePersonalized RecommendationType = "personalized"
)

type ProductRecommendation struct {
	ProductID string            `json:"product_id"`
	Title     string            `json:"title"`
	Score     float64           `json:"score"`
	Type      RecommendationType `json:"type"`
	Reason    string            `json:"reason"`
	Category  string            `json:"category"`
	Price     float64           `json:"price"`
	CreatedAt time.Time         `json:"created_at"`
}

type RecommendationContext struct {
	UserID    string             `json:"user_id"`
	Type      RecommendationType `json:"type"`
	ProductID string             `json:"product_id,omitempty"`
	SessionID string             `json:"session_id,omitempty"`
	Limit     int                `json:"limit"`
	Offset    int                `json:"offset"`
}

type Score struct {
	Collaborative float64
	ContentBased  float64
	Trending      float64
	Personalized  float64
}

type Reason string

const (
	ReasonBoughtAlsoBought Reason = "Customers who bought this also bought"
	ReasonSimilarProducts  Reason = "Similar products you might like"
	ReasonTrendingNow      Reason = "Trending now"
	ReasonPersonalized     Reason = "Personalized for you"
)
