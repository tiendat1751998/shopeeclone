package search

import "time"

type ProductDocument struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	SellerID    string    `json:"seller_id"`
	Price       float64   `json:"price"`
	Rating      float64   `json:"rating"`
	Stock       int       `json:"stock"`
	Tags        []string  `json:"tags"`
	ImageURLs   []string  `json:"image_urls"`
	CreatedAt   time.Time `json:"created_at"`
}

type SearchQuery struct {
	Query     string
	Category  string
	Brand     string
	MinPrice  float64
	MaxPrice  float64
	MinRating float64
	SortBy    string
	Page      int
	Limit     int
}

type Facet struct {
	Field  string           `json:"field"`
	Values map[string]int64 `json:"values"`
}

type SearchResult struct {
	Products []ProductDocument `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	TookMs   int64             `json:"took_ms"`
	Facets   []Facet           `json:"facets,omitempty"`
}
