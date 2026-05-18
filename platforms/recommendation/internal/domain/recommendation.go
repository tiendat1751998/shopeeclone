package domain
import "time"

type Recommendation struct { ID string `json:"id"`; UserID string `json:"user_id"`; ProductID string `json:"product_id"`; Score float64 `json:"score"`; Type string `json:"type"`; Reason string `json:"reason"`; CreatedAt time.Time `json:"created_at"` }

type RecommendationRequest struct { UserID string `json:"user_id"`; Context string `json:"context"`; ProductID string `json:"product_id,omitempty"`; Limit int `json:"limit"` }

type RecommendationResponse struct { Recommendations []Recommendation `json:"recommendations"`; TookMs int64 `json:"took_ms"` }

type UserEvent struct { UserID string `json:"user_id"`; EventType string `json:"event_type"`; ProductID string `json:"product_id"`; Timestamp time.Time `json:"timestamp"`; Metadata map[string]interface{} `json:"metadata"` }

type FeatureVector struct { UserID string `json:"user_id"`; ProductID string `json:"product_id"`; Features map[string]float64 `json:"features"`; UpdatedAt time.Time `json:"updated_at"` }

const (
	RecTypeHome       = "home"
	RecTypeProduct    = "product"
	RecTypeSimilar    = "similar"
	RecTypeTrending   = "trending"
	RecTypePersonalized = "personalized"
)

var ErrRecFailed = ErrRecommendation("recommendation_failed")
type ErrRecommendation string
func (e ErrRecommendation) Error() string { return "recommendation: " + string(e) }
