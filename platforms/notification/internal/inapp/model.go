package inapp

import "time"

type Category string

const (
	CategoryOrder    Category = "order"
	CategoryPayment  Category = "payment"
	CategoryPromo    Category = "promotion"
	CategorySystem   Category = "system"
	CategorySocial   Category = "social"
)

type Action struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	URL   string `json:"url,omitempty"`
}

type InAppNotification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Category  Category  `json:"category"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	ImageURL  string    `json:"image_url,omitempty"`
	Action    *Action   `json:"action,omitempty"`
	Read      bool      `json:"read"`
	Dismissed bool      `json:"dismissed"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

type SendInAppRequest struct {
	UserID   string   `json:"user_id"`
	Category Category `json:"category"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	ImageURL string   `json:"image_url,omitempty"`
	Action   *Action  `json:"action,omitempty"`
}
