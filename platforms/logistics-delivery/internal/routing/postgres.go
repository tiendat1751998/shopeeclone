package routing

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateRoute(ctx context.Context, rt *Route) error {
	query := `INSERT INTO routes (
		id, shipment_id, route_type, origin_id, destination_id,
		distance_km, estimated_duration_min, priority, is_active, created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, query,
		rt.ID, rt.ShipmentID, rt.RouteType, rt.OriginID, rt.DestinationID,
		rt.DistanceKm, rt.EstimatedDurationMin, rt.Priority, rt.IsActive,
		rt.CreatedAt, rt.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetRoutesByShipment(ctx context.Context, shipmentID string) ([]*Route, error) {
	rows, err := r.pool.Query(ctx, `SELECT
		id, shipment_id, route_type, origin_id, destination_id,
		distance_km, estimated_duration_min, priority, is_active, created_at, updated_at
		FROM routes WHERE shipment_id=$1 ORDER BY priority ASC`, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*Route
	for rows.Next() {
		rt := &Route{}
		if err := rows.Scan(&rt.ID, &rt.ShipmentID, &rt.RouteType, &rt.OriginID, &rt.DestinationID,
			&rt.DistanceKm, &rt.EstimatedDurationMin, &rt.Priority, &rt.IsActive,
			&rt.CreatedAt, &rt.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, rt)
	}
	return result, nil
}

func (r *PostgresRepository) GetActiveRoute(ctx context.Context, shipmentID string) (*Route, error) {
	rt := &Route{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, route_type, origin_id, destination_id,
		distance_km, estimated_duration_min, priority, is_active, created_at, updated_at
		FROM routes WHERE shipment_id=$1 AND is_active=true ORDER BY priority ASC LIMIT 1`, shipmentID).Scan(
		&rt.ID, &rt.ShipmentID, &rt.RouteType, &rt.OriginID, &rt.DestinationID,
		&rt.DistanceKm, &rt.EstimatedDurationMin, &rt.Priority, &rt.IsActive,
		&rt.CreatedAt, &rt.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *PostgresRepository) UpdateRoute(ctx context.Context, rt *Route) error {
	_, err := r.pool.Exec(ctx, `UPDATE routes SET
		route_type=$1, origin_id=$2, destination_id=$3, distance_km=$4,
		estimated_duration_min=$5, priority=$6, is_active=$7, updated_at=NOW()
		WHERE id=$8`, rt.RouteType, rt.OriginID, rt.DestinationID, rt.DistanceKm,
		rt.EstimatedDurationMin, rt.Priority, rt.IsActive, rt.ID)
	return err
}

func (r *PostgresRepository) GetZones(ctx context.Context) ([]*Zone, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, name, city, state, is_active FROM zones WHERE is_active=true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*Zone
	for rows.Next() {
		z := &Zone{}
		if err := rows.Scan(&z.ID, &z.Name, &z.City, &z.State, &z.IsActive); err != nil {
			return nil, err
		}
		result = append(result, z)
	}
	return result, nil
}

func (r *PostgresRepository) GetZoneByCity(ctx context.Context, city string) (*Zone, error) {
	z := &Zone{}
	err := r.pool.QueryRow(ctx, "SELECT id, name, city, state, is_active FROM zones WHERE city=$1 AND is_active=true LIMIT 1", city).
		Scan(&z.ID, &z.Name, &z.City, &z.State, &z.IsActive)
	if err != nil {
		return nil, err
	}
	return z, nil
}
