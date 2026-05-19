package events

import "time"

type EventType string

const (
	EventNotificationSent  EventType = "notification.sent"
	EventPushFailed        EventType = "push.failed"
	EventEmailBounced      EventType = "email.bounced"
	EventEmailOpened       EventType = "email.opened"
	EventSMSSent           EventType = "sms.sent"
	EventInAppSent         EventType = "inapp.sent"
	EventNotificationRead  EventType = "notification.read"
)

type NotificationSentEvent struct {
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	Channel        string    `json:"channel"`
	Type           string    `json:"type"`
	SentAt         time.Time `json:"sent_at"`
}

type PushFailedEvent struct {
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	DeviceToken    string    `json:"device_token"`
	Platform       string    `json:"platform"`
	Error          string    `json:"error"`
	FailedAt       time.Time `json:"failed_at"`
}

type EmailBouncedEvent struct {
	EmailID  string    `json:"email_id"`
	To       string    `json:"to"`
	Reason   string    `json:"reason"`
	BouncedAt time.Time `json:"bounced_at"`
}

type EmailOpenedEvent struct {
	EmailID  string    `json:"email_id"`
	To       string    `json:"to"`
	OpenedAt time.Time `json:"opened_at"`
}

type SMSSentEvent struct {
	SMSID  string    `json:"sms_id"`
	To     string    `json:"to"`
	SentAt time.Time `json:"sent_at"`
}

type InAppSentEvent struct {
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	Category       string    `json:"category"`
	SentAt         time.Time `json:"sent_at"`
}

type NotificationReadEvent struct {
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	ReadAt         time.Time `json:"read_at"`
}
