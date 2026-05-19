package collaborative

type UserSimilarity struct {
	UserID     string  `json:"user_id"`
	TargetID   string  `json:"target_id"`
	Similarity float64 `json:"similarity"`
}

type ItemSimilarity struct {
	ItemID     string  `json:"item_id"`
	TargetID   string  `json:"target_id"`
	Similarity float64 `json:"similarity"`
}

type RatingMatrix struct {
	UserItemRatings map[string]map[string]float64
	ItemUserRatings map[string]map[string]float64
}
