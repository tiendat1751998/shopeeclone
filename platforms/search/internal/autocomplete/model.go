package autocomplete

type Suggestion struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
	Type  string  `json:"type"`
}

type TrendQuery struct {
	Query string  `json:"query"`
	Score float64 `json:"score"`
}

type AutocompleteResult struct {
	Suggestions []Suggestion `json:"suggestions"`
	TookMs      int64        `json:"took_ms"`
}
