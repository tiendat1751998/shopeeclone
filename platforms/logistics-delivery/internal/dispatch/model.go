package dispatch

import "time"

type DispatchStatus string

const (
	DispatchPending   DispatchStatus = "pending"
	DispatchAssigned  DispatchStatus = "assigned"
	DispatchAccepted  DispatchStatus = "accepted"
	DispatchEnRoute   DispatchStatus = "en_route"
	DispatchArrived   DispatchStatus = "arrived"
	DispatchCompleted DispatchStatus = "completed"
	DispatchCancelled DispatchStatus = "cancelled"
	DispatchFailed    DispatchStatus = "failed"
)

type Dispatch struct {
	ID          string         `json:"id"`
	ShipmentID  string         `json:"shipment_id"`
	CourierID   string         `json:"courier_id"`
	ZoneID      string         `json:"zone_id"`
	Status      DispatchStatus `json:"status"`
	PickupTime  *time.Time     `json:"pickup_time,omitempty"`
	DispatchTime *time.Time    `json:"dispatch_time,omitempty"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	Notes       string         `json:"notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type DispatchFilter struct {
	Status    DispatchStatus
	CourierID string
	ZoneID    string
	Offset    int
	Limit     int
}
