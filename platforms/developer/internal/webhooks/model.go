package webhooks

import "time"

type Webhook struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	Secret         string    `json:"secret"`
	Events         []string  `json:"events"`
	IsActive       bool      `json:"is_active"`
	RetryCount     int       `json:"retry_count"`
	TimeoutSeconds int       `json:"timeout_seconds"`
	CreatedAt      time.Time `json:"created_at"`
}

type DeliveryStatus string

const (
	DeliveryPending   DeliveryStatus = "pending"
	DeliveryDelivered DeliveryStatus = "delivered"
	DeliveryFailed    DeliveryStatus = "failed"
)

type Delivery struct {
	ID           string         `json:"id"`
	WebhookID    string         `json:"webhook_id"`
	Event        string         `json:"event"`
	Status       DeliveryStatus `json:"status"`
	ResponseCode int            `json:"response_code"`
	AttemptedAt  time.Time      `json:"attempted_at"`
	RetryCount   int            `json:"retry_count"`
}
