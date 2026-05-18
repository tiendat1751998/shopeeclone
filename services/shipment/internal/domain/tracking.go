package domain

import (
	"time"

	"github.com/google/uuid"
)

type TrackingEvent struct {
	ID          string    `db:"id" json:"id"`
	ShipmentID  string    `db:"shipment_id" json:"shipment_id"`
	Status      string    `db:"status" json:"status"`
	Location    string    `db:"location" json:"location"`
	Description string    `db:"description" json:"description"`
	Timestamp   time.Time `db:"timestamp" json:"timestamp"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func NewTrackingEvent(shipmentID, status, location, description string) *TrackingEvent {
	now := time.Now().UTC()
	return &TrackingEvent{
		ID: uuid.New().String(), ShipmentID: shipmentID, Status: status,
		Location: location, Description: description, Timestamp: now, CreatedAt: now,
	}
}
