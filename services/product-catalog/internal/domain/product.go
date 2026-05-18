package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string    `db:"id" json:"id"`
	ShopID      string    `db:"shop_id" json:"shop_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description,omitempty"`
	CategoryID  string    `db:"category_id" json:"category_id"`
	Status      string    `db:"status" json:"status"`
	Currency    string    `db:"currency" json:"currency"`
	Version     int64     `db:"version" json:"version"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
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

type SKU struct {
	ID         string    `db:"id" json:"id"`
	ProductID  string    `db:"product_id" json:"product_id"`
	SKUCode    string    `db:"sku_code" json:"sku_code"`
	Attributes string    `db:"attributes" json:"attributes"`
	Price      int64     `db:"price" json:"price"`
	SalePrice  int64     `db:"sale_price" json:"sale_price,omitempty"`
	Stock      int64     `db:"stock" json:"stock"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

const (
	SKUStatusActive   = "active"
	SKUStatusInactive = "inactive"
	SKUStatusOutOfStock = "out_of_stock"
)

func NewSKU(productID, skuCode, attributes string, price int64) *SKU {
	now := time.Now()
	return &SKU{
		ID: uuid.New().String(), ProductID: productID, SKUCode: skuCode,
		Attributes: attributes, Price: price, Status: SKUStatusActive,
		CreatedAt: now, UpdatedAt: now,
	}
}

type Category struct {
	ID        string    `db:"id" json:"id"`
	ParentID  string    `db:"parent_id" json:"parent_id,omitempty"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	Level     int       `db:"level" json:"level"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	Metadata  string    `db:"metadata" json:"metadata,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func NewCategory(parentID, name, slug string, level, sortOrder int) *Category {
	now := time.Now()
	return &Category{
		ID: uuid.New().String(), ParentID: parentID, Name: name,
		Slug: slug, Level: level, SortOrder: sortOrder,
		IsActive: true, CreatedAt: now, UpdatedAt: now,
	}
}

type Attribute struct {
	ID          string `db:"id" json:"id"`
	CategoryID  string `db:"category_id" json:"category_id"`
	Name        string `db:"name" json:"name"`
	DisplayName string `db:"display_name" json:"display_name"`
	AttrType    string `db:"type" json:"type"`
	Required    bool   `db:"required" json:"required"`
	Options     string `db:"options" json:"options,omitempty"`
	SortOrder   int    `db:"sort_order" json:"sort_order"`
	IsActive    bool   `db:"is_active" json:"is_active"`
}

const (
	AttributeTypeText    = "text"
	AttributeTypeNumber  = "number"
	AttributeTypeSelect  = "select"
	AttributeTypeMulti   = "multi_select"
	AttributeTypeBoolean = "boolean"
)

type ProductMedia struct {
	ID         string    `db:"id" json:"id"`
	ProductID  string    `db:"product_id" json:"product_id"`
	MediaType  string    `db:"media_type" json:"media_type"`
	URL        string    `db:"url" json:"url"`
	Thumbnail  string    `db:"thumbnail" json:"thumbnail,omitempty"`
	SortOrder  int       `db:"sort_order" json:"sort_order"`
	IsPrimary  bool      `db:"is_primary" json:"is_primary"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

const (
	MediaTypeImage = "image"
	MediaTypeVideo = "video"
)

var (
	ErrProductNotFound  = ErrCatalog("product_not_found")
	ErrSKUNotFound      = ErrCatalog("sku_not_found")
	ErrCategoryNotFound = ErrCatalog("category_not_found")
	ErrInvalidState     = ErrCatalog("invalid_state")
	ErrDuplicateSKU     = ErrCatalog("duplicate_sku")
)

type ErrCatalog string
func (e ErrCatalog) Error() string { return "catalog: " + string(e) }
