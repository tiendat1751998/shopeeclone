package audience

import "time"

type Criteria struct {
	AgeRange       *AgeRange `json:"age_range,omitempty"`
	Gender         *string   `json:"gender,omitempty"`
	Location       *string   `json:"location,omitempty"`
	PurchaseHistory *PurchaseHistory `json:"purchase_history,omitempty"`
	LastActive     *string   `json:"last_active,omitempty"`
	Tags           []string  `json:"tags,omitempty"`
}

type AgeRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type PurchaseHistory struct {
	MinOrders   *int     `json:"min_orders,omitempty"`
	MinSpent    *float64 `json:"min_spent,omitempty"`
	CategoryIDs []string `json:"category_ids,omitempty"`
}

type Segment struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Criteria       Criteria  `json:"criteria"`
	EstimatedCount int       `json:"estimated_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserProfile struct {
	ID         string            `json:"id"`
	SegmentIDs []string          `json:"segment_ids"`
	Attributes map[string]string `json:"attributes"`
	Tags       []string          `json:"tags"`
	CreatedAt  time.Time         `json:"created_at"`
}

type CreateSegmentRequest struct {
	Name     string   `json:"name"`
	Criteria Criteria `json:"criteria"`
}
