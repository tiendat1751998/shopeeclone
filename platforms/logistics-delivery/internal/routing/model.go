package routing

import "time"

type RouteType string

const (
	RouteTypeWarehouse   RouteType = "warehouse"
	RouteTypeDispatch    RouteType = "dispatch"
	RouteTypeLastMile    RouteType = "last_mile"
	RouteTypeReturn      RouteType = "return"
)

type Route struct {
	ID              string    `json:"id"`
	ShipmentID      string    `json:"shipment_id"`
	RouteType       RouteType `json:"route_type"`
	OriginID        string    `json:"origin_id"`
	DestinationID   string    `json:"destination_id"`
	DistanceKm      float64   `json:"distance_km"`
	EstimatedDurationMin int   `json:"estimated_duration_min"`
	Priority        int       `json:"priority"`
	Waypoints       []Waypoint `json:"waypoints"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Waypoint struct {
	ID        string    `json:"id"`
	RouteID   string    `json:"route_id"`
	Sequence  int       `json:"sequence"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Action    string    `json:"action"`
	ArrivedAt *time.Time `json:"arrived_at,omitempty"`
	DepartedAt *time.Time `json:"departed_at,omitempty"`
}

type Zone struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	IsActive  bool    `json:"is_active"`
}

type RoutingAssignment struct {
	ShipmentID string `json:"shipment_id"`
	ZoneID     string `json:"zone_id"`
	WarehouseID string `json:"warehouse_id"`
	CourierID  string `json:"courier_id"`
	RouteID    string `json:"route_id"`
}
