package sdk

import "time"

type SDK struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Language         string    `json:"language"`
	Version          string    `json:"version"`
	RepositoryURL    string    `json:"repository_url"`
	DocumentationURL string    `json:"documentation_url"`
	Compatibility    string    `json:"compatibility"`
	IsLatest         bool      `json:"is_latest"`
	CreatedAt        time.Time `json:"created_at"`
}
