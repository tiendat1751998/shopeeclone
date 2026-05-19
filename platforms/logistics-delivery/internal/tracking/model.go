package tracking

import "time"

type TrackingEventType string

const (
	EventPickupScheduled   TrackingEventType = "pickup_scheduled"
	EventPickedUp          TrackingEventType = "picked_up"
	EventArrivedAtHub      TrackingEventType = "arrived_at_hub"
	EventDepartedHub       TrackingEventType = "departed_hub"
	EventInTransit         TrackingEventType = "in_transit"
	EventOutForDelivery    TrackingEventType = "out_for_delivery"
	EventDeliveryAttempted TrackingEventType = "delivery_attempted"
	EventDelivered         TrackingEventType = "delivered"
	EventFailedDelivery    TrackingEventType = "failed_delivery"
	EventReturned          TrackingEventType = "returned"
	EventOnHold            TrackingEventType = "on_hold"
	EventException         TrackingEventType = "exception"
	EventCustomsCleared    TrackingEventType = "customs_cleared"
	EventSortingComplete   TrackingEventType = "sorting_complete"
)

type TrackingEvent struct {
	ID          string            `json:"id"`
	ShipmentID  string            `json:"shipment_id"`
	EventType   TrackingEventType `json:"event_type"`
	Location    Location          `json:"location"`
	Description string            `json:"description"`
	CourierData map[string]any    `json:"courier_data,omitempty"`
	ReplayID    string            `json:"replay_id,omitempty"`
	OccurredAt  time.Time         `json:"occurred_at"`
	CreatedAt   time.Time         `json:"created_at"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}

type TrackingTimeline struct {
	ShipmentID string          `json:"shipment_id"`
	Events     []*TrackingEvent `json:"events"`
	Milestones []Milestone     `json:"milestones"`
}

type Milestone struct {
	EventType   TrackingEventType `json:"event_type"`
	AchievedAt  time.Time         `json:"achieved_at"`
	Description string            `json:"description"`
}

type TrackingFilter struct {
	ShipmentID string
	EventType  TrackingEventType
	FromDate   *time.Time
	ToDate     *time.Time
	Offset     int
	Limit      int
}
