package docs

import "time"

type DocPage struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Service   string    `json:"service"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
