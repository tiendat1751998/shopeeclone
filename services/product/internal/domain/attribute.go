package domain

import (
	"time"
)

// AttributeType defines the data type of a product attribute.
type AttributeType string

const (
	AttributeTypeText        AttributeType = "TEXT"
	AttributeTypeNumber      AttributeType = "NUMBER"
	AttributeTypeBoolean     AttributeType = "BOOLEAN"
	AttributeTypeSelect      AttributeType = "SELECT"
	AttributeTypeMultiSelect AttributeType = "MULTI_SELECT"
	AttributeTypeColor       AttributeType = "COLOR"
)

// Attribute defines a schema-level property that can be assigned to products
// within a specific category (e.g. "Material", "Screen Size").
type Attribute struct {
	ID           string        `db:"attribute_id"  json:"id"`
	CategoryID   string        `db:"category_id"   json:"category_id"`
	Name         string        `db:"name"          json:"name"`
	Type         AttributeType `db:"type"          json:"type"`
	IsRequired   bool          `db:"is_required"   json:"is_required"`
	IsFilterable bool          `db:"is_filterable" json:"is_filterable"`
	IsSearchable bool          `db:"is_searchable" json:"is_searchable"`
	SortOrder    int           `db:"sort_order"    json:"sort_order"`
	CreatedAt    time.Time     `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"    json:"updated_at"`
	Values       []AttributeValue `db:"-" json:"values,omitempty"`
}

// AttributeValue represents a predefined selectable option for an attribute.
type AttributeValue struct {
	ID           string `db:"id"             json:"id"`
	AttributeID  string `db:"attribute_id"   json:"attribute_id"`
	Value        string `db:"value"          json:"value"`
	DisplayValue string `db:"display_value"  json:"display_value,omitempty"`
	SortOrder    int    `db:"sort_order"     json:"sort_order,omitempty"`
}

// ProductAttributeValue links a product to its concrete resolved attribute values.
type ProductAttributeValue struct {
	ProductID   string `db:"product_id"   json:"product_id"`
	AttributeID string `db:"attribute_id" json:"attribute_id"`
	ValueID     string `db:"value_id"     json:"value_id,omitempty"`
	CustomValue string `db:"custom_value" json:"custom_value,omitempty"`
}

// IsPredefined returns true if the value comes from the predefined attribute values.
func (pav *ProductAttributeValue) IsPredefined() bool {
	return pav.ValueID != ""
}

// Validate checks basic attribute level constraints.
func (a *Attribute) Validate() error {
	if a.Name == "" {
		return NewDomainError("attribute name is required", "INVALID_ATTRIBUTE_NAME")
	}
	switch a.Type {
	case AttributeTypeText, AttributeTypeNumber, AttributeTypeBoolean,
		AttributeTypeSelect, AttributeTypeMultiSelect, AttributeTypeColor:
		// valid
	default:
		return NewDomainError("invalid attribute type", "INVALID_ATTRIBUTE_TYPE")
	}
	return nil
}
