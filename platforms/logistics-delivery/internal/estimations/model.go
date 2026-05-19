package estimations

import "time"

type Estimation struct {
	ID               string    `json:"id"`
	ShipmentID       string    `json:"shipment_id"`
	DistanceKm       float64   `json:"distance_km"`
	BaseDurationMin  int       `json:"base_duration_min"`
	TrafficDelayMin  int       `json:"traffic_delay_min"`
	WeatherDelayMin  int       `json:"weather_delay_min"`
	TotalDurationMin int       `json:"total_duration_min"`
	ETA              time.Time `json:"eta"`
	Confidence       float64   `json:"confidence"`
	RouteHash        string    `json:"route_hash"`
	CalculatedAt     time.Time `json:"calculated_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}

type EstimationRequest struct {
	OriginLat      float64 `json:"origin_lat"`
	OriginLng      float64 `json:"origin_lng"`
	DestLat        float64 `json:"dest_lat"`
	DestLng        float64 `json:"dest_lng"`
	PackageWeight  float64 `json:"package_weight"`
	DistanceKm     float64 `json:"distance_km"`
	TrafficFactor  float64 `json:"traffic_factor"`
}
