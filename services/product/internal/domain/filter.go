package domain

// SortDirection defines sort ordering.
type SortDirection string

const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

// SortField defines available product sort fields.
type SortField string

const (
	SortByRelevance  SortField = "relevance"
	SortByPrice      SortField = "price"
	SortByCreatedAt  SortField = "created_at"
	SortByUpdatedAt  SortField = "updated_at"
	SortBySales      SortField = "sales"
	SortByRating     SortField = "rating"
	SortByPopularity SortField = "popularity"
	SortByName       SortField = "name"
)

// ProductFilter encapsulates query parameters for searching and listing products.
type ProductFilter struct {
Page       int           `json:"page"`
	Size       int           `json:"size"`
	CategoryID string        `json:"category_id,omitempty"`
	SellerID   string        `json:"seller_id,omitempty"`
	BrandID    string        `json:"brand_id,omitempty"`
	MinPrice   float64       `json:"min_price,omitempty"`
	MaxPrice   float64       `json:"max_price,omitempty"`
	SortBy     SortField     `json:"sort_by"`
	SortOrder  SortDirection `json:"sort_order"`
	Status     string        `json:"status,omitempty"`
	Search     string        `json:"search,omitempty"`
}

// ProductList holds a paginated list of products.
type ProductList struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	Size     int       `json:"size"`
}

// Normalize ensures sane defaults for pagination and sorting.
func (f *ProductFilter) Normalize() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Size < 1 {
		f.Size = 20
	}
	if f.Size > 100 {
		f.Size = 100
	}
	if f.SortBy == "" {
		f.SortBy = SortByCreatedAt
	}
	if f.SortOrder == "" {
		f.SortOrder = SortDesc
	}
}

// Offset computes the database OFFSET for pagination.
func (f *ProductFilter) Offset() int {
	return (f.Page - 1) * f.Size
}

// TotalPages computes the total number of pages.
func (pl *ProductList) TotalPages() int {
	if pl.Size == 0 {
		return 0
	}
	pages := int(pl.Total) / pl.Size
	if int(pl.Total)%pl.Size != 0 {
		pages++
	}
	return pages
}

// HasNext indicates whether more pages exist.
func (pl *ProductList) HasNext() bool {
	return pl.Page < pl.TotalPages()
}

// HasPrevious indicates whether a previous page exists.
func (pl *ProductList) HasPrevious() bool {
	return pl.Page > 1
}

// IsEmpty returns true if no products were returned.
func (pl *ProductList) IsEmpty() bool {
	return len(pl.Products) == 0
}
