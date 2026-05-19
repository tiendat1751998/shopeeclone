package pickups

import "time"

type PickupStatus string

const (
	PickupScheduled  PickupStatus = "scheduled"
	PickupAssigned   PickupStatus = "assigned"
	PickupInProgress PickupStatus = "in_progress"
	PickupCompleted  PickupStatus = "completed"
	PickupFailed     PickupStatus = "failed"
	PickupCancelled  PickupStatus = "cancelled"
)

type Pickup struct {
	ID          string       `json:"id"`
	ShipmentID  string       `json:"shipment_id"`
	FulfillmentID string     `json:"fulfillment_id"`
	CourierID   string       `json:"courier_id"`
	Status      PickupStatus `json:"status"`
	Address     string       `json:"address"`
	Latitude    float64      `json:"latitude"`
	Longitude   float64      `json:"longitude"`
	ScheduledAt *time.Time   `json:"scheduled_at,omitempty"`
	PickedUpAt  *time.Time   `json:"picked_up_at,omitempty"`
	Notes       string       `json:"notes,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
