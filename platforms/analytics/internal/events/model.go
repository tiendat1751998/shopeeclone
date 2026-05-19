package events

import "time"

type EventType string

const (
	EventPageview   EventType = "pageview"
	EventPurchase   EventType = "purchase"
	EventAddToCart  EventType = "add_to_cart"
	EventLogin      EventType = "login"
	EventSignup     EventType = "signup"
	EventSearch     EventType = "search"
	EventClick      EventType = "click"
	EventCheckout   EventType = "checkout"
	EventLogout     EventType = "logout"
	EventViewItem   EventType = "view_item"
)

type EventContext struct {
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Device    string `json:"device,omitempty"`
	Referrer  string `json:"referrer,omitempty"`
	OS        string `json:"os,omitempty"`
	Browser   string `json:"browser,omitempty"`
}

type AnalyticsEvent struct {
	EventID    string                 `json:"event_id"`
	EventType  EventType              `json:"event_type"`
	UserID     string                 `json:"user_id"`
	SessionID  string                 `json:"session_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Context    EventContext           `json:"context,omitempty"`
	Source     string                 `json:"source,omitempty"`
	Country    string                 `json:"country,omitempty"`
	Device     string                 `json:"device,omitempty"`
	Campaign   string                 `json:"campaign,omitempty"`
	Revenue    float64                `json:"revenue,omitempty"`
}

type BatchIngestRequest struct {
	Events []AnalyticsEvent `json:"events"`
}
