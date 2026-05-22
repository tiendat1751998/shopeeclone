package http

type ProductResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	CategoryID  string             `json:"category_id"`
	Brand       string             `json:"brand"`
	Status      string             `json:"status"`
	Condition   string             `json:"condition"`
	Weight      float64            `json:"weight"`
	Dimensions  string             `json:"dimensions"`
	Version     int64              `json:"version"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	ShopID      string             `json:"shop_id"`
	SKUs        []SKUResponse      `json:"skus,omitempty"`
	Media       []MediaResponse    `json:"media,omitempty"`
	Attributes  map[string]string  `json:"attributes,omitempty"`
	SoldCount   int64              `json:"sold_count"`
}

type SKUResponse struct {
	ID           string              `json:"id"`
	ProductID    string              `json:"product_id"`
	Name         string              `json:"name"`
	Price        float64             `json:"price"`
	ComparePrice float64             `json:"compare_price"`
	Currency     string              `json:"currency"`
	Stock        int32               `json:"stock"`
	ReservedStock int32              `json:"reserved_stock"`
	Weight       float64             `json:"weight"`
	Dimensions   string              `json:"dimensions"`
	Status       string              `json:"status"`
	SortOrder    int32               `json:"sort_order"`
	Attributes   map[string]string   `json:"attributes,omitempty"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
}

type MediaResponse struct {
	ID            string `json:"id"`
	ProductID     string `json:"product_id"`
	Type          string `json:"type"`
	URL           string `json:"url"`
	ThumbnailURL  string `json:"thumbnail_url,omitempty"`
	AltText       string `json:"alt_text"`
	SortOrder     int32  `json:"sort_order"`
	Status        string `json:"status"`
}

type CategoryResponse struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Slug         string              `json:"slug"`
	ParentID     string              `json:"parent_id,omitempty"`
	Description  string              `json:"description"`
	ImageURL     string              `json:"image_url"`
	SortOrder    int32               `json:"sort_order"`
	IsActive     bool                `json:"is_active"`
	Depth        int32               `json:"depth"`
	Path         string              `json:"path"`
	Children     []CategoryResponse  `json:"children,omitempty"`
	ProductCount int64               `json:"product_count"`
}

type ProductListResponse struct {
	Products    []ProductResponse `json:"products"`
	Total       int64             `json:"total"`
	Page        int               `json:"page"`
	Size        int               `json:"size"`
}

type SearchResultResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
