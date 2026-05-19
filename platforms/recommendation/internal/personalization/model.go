package personalization

type UserProfile struct {
	UserID            string             `json:"user_id"`
	CategoryWeights   map[string]float64 `json:"category_weights"`
	PreferredBrands   map[string]float64 `json:"preferred_brands"`
	PriceRangeMin     float64            `json:"price_range_min"`
	PriceRangeMax     float64            `json:"price_range_max"`
	PreferredPriceMid float64            `json:"preferred_price_mid"`
	InterestVector    map[string]float64 `json:"interest_vector"`
	TotalInteractions int                `json:"total_interactions"`
}

type UserPreference struct {
	Category string  `json:"category"`
	Weight   float64 `json:"weight"`
	Brand    string  `json:"brand"`
	MaxPrice float64 `json:"max_price"`
}

type InterestVector struct {
	Tags   map[string]float64 `json:"tags"`
	Categories []string        `json:"categories"`
}
