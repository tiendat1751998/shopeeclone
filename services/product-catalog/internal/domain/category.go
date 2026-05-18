package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          string    `db:"id" json:"id"`
	ParentID    *string   `db:"parent_id" json:"parent_id,omitempty"`
	Name        string    `db:"name" json:"name"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	ImageURL    string    `db:"image_url" json:"image_url"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	Depth       int       `db:"depth" json:"depth"`
	Path        string    `db:"path" json:"path"`
	Version     int       `db:"version" json:"version"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func NewCategory(name, slug, description string, parentID *string, depth int) *Category {
	now := time.Now().UTC()
	path := slug
	if parentID != nil {
		path = *parentID + "/" + slug
	}
	return &Category{
		ID: uuid.New().String(), ParentID: parentID, Name: name, Slug: slug,
		Description: description, IsActive: true, Depth: depth, Path: path,
		Version: 1, CreatedAt: now, UpdatedAt: now,
	}
}

func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

func (c *Category) UpdatePath(parentPath string) {
	if parentPath != "" {
		c.Path = parentPath + "/" + c.Slug
	} else {
		c.Path = c.Slug
	}
	c.UpdatedAt = time.Now().UTC()
}
