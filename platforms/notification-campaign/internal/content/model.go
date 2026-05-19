package content

import "time"

type ContentTemplate struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Channel      string    `json:"channel"`
	Subject      string    `json:"subject"`
	Body         string    `json:"body"`
	Variables    []string  `json:"variables"`
	PreviewText  string    `json:"preview_text,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Variant struct {
	ID               string             `json:"id"`
	TemplateID       string             `json:"template_id"`
	Name             string             `json:"name"`
	Modifications    map[string]string  `json:"modifications"`
	TrafficPercentage int                `json:"traffic_percentage"`
	CreatedAt        time.Time          `json:"created_at"`
}

type CreateTemplateRequest struct {
	Name        string   `json:"name"`
	Channel     string   `json:"channel"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Variables   []string `json:"variables"`
	PreviewText string   `json:"preview_text,omitempty"`
}

type CreateVariantRequest struct {
	TemplateID        string            `json:"template_id"`
	Name              string            `json:"name"`
	Modifications     map[string]string `json:"modifications"`
	TrafficPercentage int               `json:"traffic_percentage"`
}

type RenderRequest struct {
	TemplateID string                 `json:"template_id"`
	Variables  map[string]interface{} `json:"variables"`
}
