package domain

import "encoding/json"

type CategoryCreatedEvent struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	ParentID   string `json:"parent_id,omitempty"`
}

func NewCategoryCreatedEvent(c *Category) *CategoryCreatedEvent {
	e := &CategoryCreatedEvent{
		CategoryID: c.CategoryID,
		Name:       c.Name,
	}
	if c.ParentID != "" {
		e.ParentID = c.ParentID
	}
	return e
}

func (e *CategoryCreatedEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}
