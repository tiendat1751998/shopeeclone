package routing

import (
	"context"
	"math"
	"time"
)

type Repository interface {
	CreateRoute(ctx context.Context, r *Route) error
	GetRoutesByShipment(ctx context.Context, shipmentID string) ([]*Route, error)
	GetActiveRoute(ctx context.Context, shipmentID string) (*Route, error)
	UpdateRoute(ctx context.Context, r *Route) error
	GetZones(ctx context.Context) ([]*Zone, error)
	GetZoneByCity(ctx context.Context, city string) (*Zone, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AssignWarehouse(ctx context.Context, shipmentID, city, state string) (*RoutingAssignment, error) {
	zone, err := s.repo.GetZoneByCity(ctx, city)
	if err != nil {
		zone = &Zone{
			ID:   generateZoneID(city, state),
			Name: city + " Zone",
			City: city,
			State: state,
		}
	}
	return &RoutingAssignment{
		ShipmentID:  shipmentID,
		ZoneID:      zone.ID,
		WarehouseID: zone.ID + "-wh-001",
	}, nil
}

func (s *Service) CreateRoute(ctx context.Context, route *Route) error {
	if len(route.Waypoints) > 1 {
		route.EstimatedDurationMin = calculateDuration(route.Waypoints)
	}
	route.CreatedAt = time.Now().UTC()
	route.UpdatedAt = route.CreatedAt
	return s.repo.CreateRoute(ctx, route)
}

func (s *Service) GetRoutesByShipment(ctx context.Context, shipmentID string) ([]*Route, error) {
	return s.repo.GetRoutesByShipment(ctx, shipmentID)
}

func (s *Service) OptimizeWaypoints(ctx context.Context, waypoints []Waypoint) ([]Waypoint, error) {
	if len(waypoints) < 3 {
		return waypoints, nil
	}
	sorted := make([]Waypoint, len(waypoints))
	copy(sorted, waypoints)
	for i := 1; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted)-1; j++ {
			distI := haversine(sorted[i-1].Latitude, sorted[i-1].Longitude, sorted[i].Latitude, sorted[i].Longitude) +
				haversine(sorted[i].Latitude, sorted[i].Longitude, sorted[i+1].Latitude, sorted[i+1].Longitude)
			distJ := haversine(sorted[i-1].Latitude, sorted[i-1].Longitude, sorted[j].Latitude, sorted[j].Longitude) +
				haversine(sorted[j].Latitude, sorted[j].Longitude, sorted[i+1].Latitude, sorted[i+1].Longitude)
			if distJ < distI {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	for i := range sorted {
		sorted[i].Sequence = i
	}
	return sorted, nil
}

func generateZoneID(city, state string) string {
	return city + "-" + state
}

func calculateDuration(waypoints []Waypoint) int {
	if len(waypoints) < 2 {
		return 0
	}
	var totalDist float64
	for i := 1; i < len(waypoints); i++ {
		totalDist += haversine(waypoints[i-1].Latitude, waypoints[i-1].Longitude, waypoints[i].Latitude, waypoints[i].Longitude)
	}
	return int(totalDist / 40.0 * 60.0)
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
