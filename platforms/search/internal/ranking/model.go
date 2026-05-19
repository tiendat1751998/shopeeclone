package ranking

type RankScore float64

type RankingFactor struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	Value  float64 `json:"value"`
}

type RankingConfig struct {
	TitleBoost      float64 `json:"title_boost"`
	CategoryBoost   float64 `json:"category_boost"`
	RatingBoost     float64 `json:"rating_boost"`
	RecencyBoost    float64 `json:"recency_boost"`
	PopularityBoost float64 `json:"popularity_boost"`
	ClickBoost      float64 `json:"click_boost"`
}

type ClickSignal struct {
	ProductID string  `json:"product_id"`
	Query     string  `json:"query"`
	CTR       float64 `json:"ctr"`
	Count     int64   `json:"count"`
}

func DefaultRankingConfig() RankingConfig {
	return RankingConfig{
		TitleBoost:      3.0,
		CategoryBoost:   2.0,
		RatingBoost:     1.5,
		RecencyBoost:    1.0,
		PopularityBoost: 1.0,
		ClickBoost:      2.0,
	}
}
