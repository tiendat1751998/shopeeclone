package core

import "time"

type EventType string

const (
	EventLogin    EventType = "login"
	EventPayment  EventType = "payment"
	EventOrder    EventType = "order"
	EventRegister EventType = "registration"
)

type FraudEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	UserID    string                 `json:"user_id"`
	IP        string                 `json:"ip"`
	DeviceID  string                 `json:"device_id"`
	Amount    float64                `json:"amount"`
	Currency  string                 `json:"currency"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}
