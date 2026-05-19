package notifier

import "time"

type NotificationType string

const (
	TypeOrderConfirmation NotificationType = "order_confirmation"
	TypeShipmentUpdate    NotificationType = "shipment_update"
	TypePaymentReceipt    NotificationType = "payment_receipt"
	TypePromotion         NotificationType = "promotion"
	TypePriceDrop         NotificationType = "price_drop"
	TypeBackInStock       NotificationType = "back_in_stock"
	TypeWelcome           NotificationType = "welcome"
	TypeResetPassword     NotificationType = "reset_password"
	TypeVerification      NotificationType = "verification"
	TypeAdminAlert        NotificationType = "admin_alert"
)

type Channel string

const (
	ChannelPush  Channel = "push"
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelInApp Channel = "inapp"
)

type DeliveryStatus string

const (
	StatusPending   DeliveryStatus = "pending"
	StatusSent      DeliveryStatus = "sent"
	StatusDelivered DeliveryStatus = "delivered"
	StatusFailed    DeliveryStatus = "failed"
	StatusBounced   DeliveryStatus = "bounced"
	StatusRead      DeliveryStatus = "read"
)

type Priority int

const (
	PriorityLow    Priority = 0
	PriorityNormal Priority = 1
	PriorityHigh   Priority = 2
	PriorityUrgent Priority = 3
)

type Notification struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	Type      NotificationType `json:"type"`
	Channel   Channel          `json:"channel"`
	Title     string           `json:"title"`
	Body      string           `json:"body"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Status    DeliveryStatus   `json:"status"`
	Priority  Priority         `json:"priority"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type SendNotificationRequest struct {
	UserID   string                 `json:"user_id"`
	Type     NotificationType       `json:"type"`
	Channel  Channel                `json:"channel"`
	Title    string                 `json:"title"`
	Body     string                 `json:"body"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Priority Priority               `json:"priority,omitempty"`
	TemplateID string               `json:"template_id,omitempty"`
}

type Template struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      NotificationType  `json:"type"`
	Subject   string            `json:"subject"`
	Body      string            `json:"body"`
	Variables []string          `json:"variables"`
	Version   int               `json:"version"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
