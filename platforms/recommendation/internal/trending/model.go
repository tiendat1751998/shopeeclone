package trending

import "time"

type TrendingScore struct {
	ProductID string  `json:"product_id"`
	Score     float64 `json:"score"`
	Velocity  float64 `json:"velocity"`
}

type TrendWindow struct {
	WindowSize time.Duration `json:"window_size"`
	Weight     float64       `json:"weight"`
}
