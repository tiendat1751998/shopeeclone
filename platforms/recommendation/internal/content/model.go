package content

type ProductFeatures struct {
	ProductID  string   `json:"product_id"`
	Category   string   `json:"category"`
	ParentCategory string `json:"parent_category"`
	Tags       []string `json:"tags"`
	Price      float64  `json:"price"`
	Brand      string   `json:"brand"`
}

type CategoryEmbedding struct {
	Category     string  `json:"category"`
	Parent       string  `json:"parent"`
	Level        int     `json:"level"`
}

type TagVector struct {
	Tags map[string]float64 `json:"tags"`
}
