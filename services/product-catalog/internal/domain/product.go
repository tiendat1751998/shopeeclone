package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string           `db:"id" json:"id"`
	ShopID      string           `db:"shop_id" json:"shop_id"`
	Name        string           `db:"name" json:"name"`
	Description string           `db:"description" json:"description,omitempty"`
	CategoryID  string           `db:"category_id" json:"category_id"`
	Brand       string           `db:"brand" json:"brand,omitempty"`
	IdempotencyKey string         `db:"idempotency_key" json:"idempotency_key,omitempty"`
	Status      string           `db:"status" json:"status"`
	Condition   string           `db:"condition" json:"condition,omitempty"`
	Weight      float64          `db:"weight" json:"weight,omitempty"`
	Dimensions  string           `db:"dimensions" json:"dimensions,omitempty"`
	Metadata    json.RawMessage  `db:"metadata" json:"metadata,omitempty"`
	Currency    string           `db:"currency" json:"currency"`
	Version     int64            `db:"version" json:"version"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time       `db:"deleted_at" json:"deleted_at,omitempty"`
}

const (
	ProductStatusDraft     = "draft"
	ProductStatusActive    = "active"
	ProductStatusInactive  = "inactive"
	ProductStatusArchived  = "archived"
	ProductStatusModerated = "moderated"
)

func NewProduct(shopID, name, description, categoryID, currency string) *Product {
	now := time.Now()
	return &Product{
		ID: uuid.New().String(), ShopID: shopID, Name: name,
		Description: description, CategoryID: categoryID,
		Status: ProductStatusDraft, Currency: currency,
		Version: 1, CreatedAt: now, UpdatedAt: now,
	}
}

func (p *Product) Activate() error {
	if p.Status != ProductStatusDraft && p.Status != ProductStatusInactive {
		return fmt.Errorf("%w: cannot activate product in status %s", ErrInvalidState, p.Status)
	}
	p.Status = ProductStatusActive
	p.Version++
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) Archive() error {
	if p.Status == ProductStatusArchived {
		return fmt.Errorf("%w: product already archived", ErrInvalidState)
	}
	p.Status = ProductStatusArchived
	p.Version++
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) Update(name, description string) {
	if name != "" { p.Name = name }
	if description != "" { p.Description = description }
	p.Version++
	p.UpdatedAt = time.Now()
}

type ProductMedia struct {
	ID        string    `db:"id" json:"id"`
	ProductID string    `db:"product_id" json:"product_id"`
	MediaType string    `db:"media_type" json:"media_type"`
	URL       string    `db:"url" json:"url"`
	Thumbnail string    `db:"thumbnail" json:"thumbnail,omitempty"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	IsPrimary bool      `db:"is_primary" json:"is_primary"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewProductMedia(productID, mediaType, url string, sortOrder int) *ProductMedia {
	return &ProductMedia{
		ID: uuid.New().String(), ProductID: productID,
		MediaType: mediaType, URL: url, SortOrder: sortOrder,
		CreatedAt: time.Now(),
	}
}

