package events

import "time"

type EventType string

const (
	EventRecommendationRequested EventType = "recommendation.requested"
	EventItemClicked             EventType = "item.clicked"
	EventItemPurchased           EventType = "item.purchased"
	EventItemViewed              EventType = "item.viewed"
)

type RecommendationRequested struct {
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
	Limit     int    `json:"limit"`
	Timestamp time.Time `json:"timestamp"`
}

type ItemClicked struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	SessionID string    `json:"session_id"`
	Timestamp time.Time `json:"timestamp"`
}

type ItemPurchased struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}
