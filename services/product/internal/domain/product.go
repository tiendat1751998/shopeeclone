package domain

import (
	"time"
)

// ProductStatus represents the lifecycle state of a product.
type ProductStatus string

const (
	ProductStatusDraft         ProductStatus = "DRAFT"
	ProductStatusPendingReview ProductStatus = "PENDING_REVIEW"
	ProductStatusActive        ProductStatus = "ACTIVE"
	ProductStatusInactive      ProductStatus = "INACTIVE"
	ProductStatusRejected      ProductStatus = "REJECTED"
	ProductStatusDeleted       ProductStatus = "DELETED"
)

// SKUStatus represents the availability state of an SKU.
type SKUStatus string

const (
	SKUStatusActive     SKUStatus = "ACTIVE"
	SKUStatusInactive   SKUStatus = "INACTIVE"
	SKUStatusOutOfStock SKUStatus = "OUT_OF_STOCK"
)

// Product is the aggregate root representing a sellable item (SPU level).
type Product struct {
	ID          int64          `db:"id"           json:"id"`
	SPUID       string         `db:"spu_id"       json:"spu_id"`
	Title       string         `db:"title"        json:"title"`
	Description string         `db:"description"  json:"description"`
	CategoryID  string         `db:"category_id"  json:"category_id"`
	BrandID     string         `db:"brand_id"     json:"brand_id"`
	SellerID    string         `db:"seller_id"    json:"seller_id"`
	Status      ProductStatus  `db:"status"       json:"status"`
	CreatedAt   time.Time      `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"   json:"updated_at"`
	DeletedAt   *time.Time     `db:"deleted_at"   json:"deleted_at,omitempty"`
	SKUs        []SKU          `db:"-"            json:"skus,omitempty"`
	Images      []ProductImage `db:"-"            json:"images,omitempty"`
}

// SKU represents a specific sellable variant of a product.
type SKU struct {
	ID        int64     `db:"id"         json:"id"`
	SPUID     string    `db:"spu_id"     json:"spu_id"`
	SKUID     string    `db:"sku_id"     json:"sku_id"`
	Price     float64   `db:"price"      json:"price"`
	SalePrice float64   `db:"sale_price" json:"sale_price"`
	Stock     int32     `db:"stock"      json:"stock"`
	Weight    float64   `db:"weight"     json:"weight"`
	Length    float64   `db:"length"     json:"length"`
	Width     float64   `db:"width"      json:"width"`
	Height    float64   `db:"height"     json:"height"`
	Status    SKUStatus `db:"status"     json:"status"`
	Version   int32     `db:"version"    json:"version"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// ProductImage represents an image associated with a product.
type ProductImage struct {
	ID        int64  `db:"id"         json:"id"`
	SPUID     string `db:"spu_id"     json:"spu_id"`
	URL       string `db:"url"        json:"url"`
	AltText   string `db:"alt_text"   json:"alt_text"`
	SortOrder int    `db:"sort_order" json:"sort_order"`
	IsPrimary bool   `db:"is_primary" json:"is_primary"`
}

// IsAvailable returns true if the product is visible and purchasable.
func (p *Product) IsAvailable() bool {
	return p.Status == ProductStatusActive && p.DeletedAt == nil
}

// IsListable returns true if the product should appear in catalog listings.
func (p *Product) IsListable() bool {
	return p.Status == ProductStatusActive
}

// HasStock returns true if at least one SKU has stock available.
func (p *Product) HasStock() bool {
	for _, sku := range p.SKUs {
		if sku.Status == SKUStatusActive && sku.Stock > 0 {
			return true
		}
	}
	return false
}

// PrimaryImage returns the primary image URL for the product.
func (p *Product) PrimaryImage() string {
	for _, img := range p.Images {
		if img.IsPrimary {
			return img.URL
		}
	}
	if len(p.Images) > 0 {
		return p.Images[0].URL
	}
	return ""
}

// EffectivePrice returns the sale price if set, otherwise the regular price.
func (s *SKU) EffectivePrice() float64 {
	if s.SalePrice > 0 {
		return s.SalePrice
	}
	return s.Price
}

// IsAvailable returns true if the SKU can be purchased.
func (s *SKU) IsAvailable() bool {
	return s.Status == SKUStatusActive && s.Stock > 0
}

// Volume returns the volumetric dimensions of the SKU.
func (s *SKU) Volume() float64 {
	return s.Length * s.Width * s.Height
}
