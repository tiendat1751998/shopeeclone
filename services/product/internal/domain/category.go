package domain

import "time"

// CategoryStatus represents the state of a category.
type CategoryStatus string

const (
	CategoryStatusActive   CategoryStatus = "active"
	CategoryStatusInactive CategoryStatus = "inactive"
)

// Category represents a hierarchical product classification node.
type Category struct {
	ID         int64     `db:"id"          json:"id"`
	CategoryID string    `db:"category_id" json:"category_id"`
	Name       string    `db:"name"        json:"name"`
	Slug       string    `db:"slug"        json:"slug"`
	ParentID   string    `db:"parent_id"   json:"parent_id,omitempty"`
	Level      int       `db:"level"       json:"level"`
	SortOrder  int       `db:"sort_order"  json:"sort_order"`
	ImageURL   string    `db:"image_url"   json:"image_url,omitempty"`
	IsActive   bool      `db:"is_active"   json:"is_active"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"  json:"updated_at"`
}

// CategoryTree provides a recursive hierarchical view of categories.
type CategoryTree struct {
	Roots    []*CategoryTreeNode `json:"roots"`
}

// CategoryTreeNode is a node in the category tree.
type CategoryTreeNode struct {
	Category Category             `json:"category"`
	Children []*CategoryTreeNode `json:"children,omitempty"`
}

// IsRoot returns true if this is a top-level category.
func (c *Category) IsRoot() bool {
	return c.ParentID == ""
}

// HasChildren returns true if the node has child categories.
func (ct *CategoryTreeNode) HasChildren() bool {
	return len(ct.Children) > 0
}

// WalkDepthFirst traverses the tree depth-first, calling fn on each node.
func (ct *CategoryTreeNode) WalkDepthFirst(fn func(node *CategoryTreeNode)) {
	if ct == nil {
		return
	}
	fn(ct)
	for _, child := range ct.Children {
		child.WalkDepthFirst(fn)
	}
}

// AllCategoryIDs collects every category ID in the subtree (inclusive).
func (ct *CategoryTreeNode) AllCategoryIDs() []string {
	var ids []string
	ct.WalkDepthFirst(func(node *CategoryTreeNode) {
		ids = append(ids, node.Category.CategoryID)
	})
	return ids
}

// FindByCategoryID searches the tree for a node by category ID.
func (ct *CategoryTreeNode) FindByCategoryID(id string) *CategoryTreeNode {
	if ct == nil {
		return nil
	}
	if ct.Category.CategoryID == id {
		return ct
	}
	for _, child := range ct.Children {
		if found := child.FindByCategoryID(id); found != nil {
			return found
		}
	}
	return nil
}
