package domain

import (
	"time"

	"github.com/google/uuid"
)

type AttributeType string

const (
	AttributeTypeText    AttributeType = "text"
	AttributeTypeNumber  AttributeType = "number"
	AttributeTypeBoolean AttributeType = "boolean"
	AttributeTypeSelect  AttributeType = "select"
	AttributeTypeMultiSelect AttributeType = "multi_select"
)

type Attribute struct {
	ID          string        `db:"id" json:"id"`
	CategoryID  string        `db:"category_id" json:"category_id"`
	Name        string        `db:"name" json:"name"`
	Slug        string        `db:"slug" json:"slug"`
	Type        AttributeType `db:"type" json:"type"`
	IsRequired  bool          `db:"is_required" json:"is_required"`
	IsFilterable bool         `db:"is_filterable" json:"is_filterable"`
	IsSearchable bool         `db:"is_searchable" json:"is_searchable"`
	Options     []byte        `db:"options" json:"options,omitempty"`
	SortOrder   int           `db:"sort_order" json:"sort_order"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at" json:"updated_at"`
}

func NewAttribute(categoryID, name, slug string, attrType AttributeType) *Attribute {
	now := time.Now().UTC()
	return &Attribute{
		ID: uuid.New().String(), CategoryID: categoryID, Name: name,
		Slug: slug, Type: attrType, CreatedAt: now, UpdatedAt: now,
	}
}

type ProductAttribute struct {
	ID          string `db:"id" json:"id"`
	ProductID   string `db:"product_id" json:"product_id"`
	AttributeID string `db:"attribute_id" json:"attribute_id"`
	Value       string `db:"value" json:"value"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func NewProductAttribute(productID, attributeID, value string) *ProductAttribute {
	return &ProductAttribute{
		ID: uuid.New().String(), ProductID: productID,
		AttributeID: attributeID, Value: value, CreatedAt: time.Now().UTC(),
	}
}
