package template

import "time"

type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Variables []string  `json:"variables"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TemplateVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
}

type TemplateVersion struct {
	ID        string    `json:"id"`
	TemplateID string   `json:"template_id"`
	Version   int       `json:"version"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTemplateRequest struct {
	Name      string   `json:"name"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
	Variables []string `json:"variables"`
}

type UpdateTemplateRequest struct {
	Subject   *string   `json:"subject,omitempty"`
	Body      *string   `json:"body,omitempty"`
	Variables *[]string `json:"variables,omitempty"`
}
