package couriers

import "time"

type CourierStatus string

const (
	CourierAvailable  CourierStatus = "available"
	CourierBusy       CourierStatus = "busy"
	CourierOffline    CourierStatus = "offline"
	CourierOnBreak    CourierStatus = "on_break"
)

type CourierProvider string

const (
	ProviderInternal CourierProvider = "internal"
	ProviderExternal CourierProvider = "external"
)

type Courier struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Phone        string          `json:"phone"`
	Provider     CourierProvider `json:"provider"`
	Status       CourierStatus   `json:"status"`
	ZoneID       string          `json:"zone_id"`
	CurrentLat   float64         `json:"current_lat"`
	CurrentLng   float64         `json:"current_lng"`
	LastSeenAt   *time.Time      `json:"last_seen_at,omitempty"`
	MaxCapacity  int             `json:"max_capacity"`
	CurrentLoad  int             `json:"current_load"`
	Rating       float64         `json:"rating"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type CourierUpdate struct {
	CourierID string         `json:"courier_id"`
	Status    CourierStatus   `json:"status,omitempty"`
	Latitude  float64         `json:"latitude,omitempty"`
	Longitude float64         `json:"longitude,omitempty"`
	Load      int             `json:"load,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

type WebhookPayload struct {
	Provider   string         `json:"provider"`
	EventType  string         `json:"event_type"`
	CourierID  string         `json:"courier_id"`
	ShipmentID string         `json:"shipment_id"`
	Data       map[string]any `json:"data"`
	Signature  string         `json:"signature"`
	ReceivedAt time.Time      `json:"received_at"`
}
