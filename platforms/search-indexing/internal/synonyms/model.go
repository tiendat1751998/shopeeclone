package synonyms

type SynonymSet struct {
	ID       string   `json:"id"`
	Words    []string `json:"words"`
	Language string   `json:"language"`
	IsActive bool     `json:"is_active"`
}

type SynonymGraph struct {
	Edges map[string][]string `json:"edges"`
}
