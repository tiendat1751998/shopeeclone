package collabvector

import "time"

type Interaction struct {
	UserID          string    `json:"user_id"`
	ItemID          string    `json:"item_id"`
	InteractionType string    `json:"interaction_type"`
	Weight          float64   `json:"weight"`
	Timestamp       time.Time `json:"timestamp"`
}

type InteractionMatrix struct {
	UserItemRatings map[string]map[string]float64
	UserIndex       map[string]int
	ItemIndex       map[string]int
	Users           []string
	Items           []string
}

type LatentFactors struct {
	UserFactors [][]float64 `json:"user_factors"`
	ItemFactors [][]float64 `json:"item_factors"`
}

type FactorRecommendation struct {
	ItemID string  `json:"item_id"`
	Score  float64 `json:"score"`
}
