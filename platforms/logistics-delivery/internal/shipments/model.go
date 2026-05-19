package shipments

import "time"

type ShipmentStatus string

const (
	StatusPending          ShipmentStatus = "pending"
	StatusConfirmed        ShipmentStatus = "confirmed"
	StatusPickedUp         ShipmentStatus = "picked_up"
	StatusInTransit        ShipmentStatus = "in_transit"
	StatusOutForDelivery   ShipmentStatus = "out_for_delivery"
	StatusDelivered        ShipmentStatus = "delivered"
	StatusFailed           ShipmentStatus = "failed"
	StatusCancelled        ShipmentStatus = "cancelled"
	StatusReturned         ShipmentStatus = "returned"
	StatusLost             ShipmentStatus = "lost"
	StatusOnHold           ShipmentStatus = "on_hold"
	StatusPartialDelivered ShipmentStatus = "partial_delivered"
	StatusAwaitingPickup   ShipmentStatus = "awaiting_pickup"
)

var validTransitions = map[ShipmentStatus][]ShipmentStatus{
	StatusPending:          {StatusConfirmed, StatusCancelled},
	StatusConfirmed:        {StatusAwaitingPickup, StatusCancelled},
	StatusAwaitingPickup:   {StatusPickedUp, StatusCancelled},
	StatusPickedUp:         {StatusInTransit, StatusFailed},
	StatusInTransit:        {StatusOutForDelivery, StatusOnHold, StatusFailed, StatusLost},
	StatusOutForDelivery:   {StatusDelivered, StatusFailed, StatusOnHold, StatusPartialDelivered},
	StatusDelivered:        {},
	StatusFailed:           {},
	StatusCancelled:        {},
	StatusReturned:         {},
	StatusLost:             {},
	StatusOnHold:           {StatusInTransit, StatusCancelled, StatusReturned},
	StatusPartialDelivered: {StatusDelivered, StatusFailed, StatusOnHold},
}

func IsValidTransition(from, to ShipmentStatus) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

type Shipment struct {
	ID                string         `json:"id"`
	OrderID           string         `json:"order_id"`
	CustomerID        string         `json:"customer_id"`
	WarehouseID       string         `json:"warehouse_id"`
	CourierID         string         `json:"courier_id,omitempty"`
	Status            ShipmentStatus `json:"status"`
	OriginAddress     Address        `json:"origin_address"`
	DestinationAddress Address       `json:"destination_address"`
	Packages          []Package      `json:"packages"`
	TotalWeight       float64        `json:"total_weight"`
	TotalVolume       float64        `json:"total_volume"`
	EstimatedDistance float64        `json:"estimated_distance"`
	EstimatedETA      *time.Time     `json:"estimated_eta,omitempty"`
	ActualDeliveredAt *time.Time     `json:"actual_delivered_at,omitempty"`
	CourierNotes      string         `json:"courier_notes,omitempty"`
	ReplayID          string         `json:"replay_id,omitempty"`
	Version           int            `json:"version"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type Address struct {
	Street     string  `json:"street"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	ZipCode    string  `json:"zip_code"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type Package struct {
	ID          string  `json:"id"`
	Weight      float64 `json:"weight"`
	Length      float64 `json:"length"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Description string  `json:"description"`
}

type ShipmentFilter struct {
	Status      ShipmentStatus
	OrderID     string
	CourierID   string
	CustomerID  string
	WarehouseID string
	Offset      int
	Limit       int
}

type StatusTransition struct {
	ShipmentID    string         `json:"shipment_id"`
	FromStatus    ShipmentStatus `json:"from_status"`
	ToStatus      ShipmentStatus `json:"to_status"`
	Reason        string         `json:"reason,omitempty"`
	ReplayID      string         `json:"replay_id,omitempty"`
	TransitionedAt time.Time     `json:"transitioned_at"`
}
