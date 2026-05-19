package estimations

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

func (r *PostgresRepository) Create(ctx context.Context, e *Estimation) error {
	query := `INSERT INTO estimations (
		id, shipment_id, distance_km, base_duration_min, traffic_delay_min,
		weather_delay_min, total_duration_min, eta, confidence, route_hash,
		calculated_at, expires_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, query,
		e.ID, e.ShipmentID, e.DistanceKm, e.BaseDurationMin, e.TrafficDelayMin,
		e.WeatherDelayMin, e.TotalDurationMin, e.ETA, e.Confidence, e.RouteHash,
		e.CalculatedAt, e.ExpiresAt)
	return err
}

func (r *PostgresRepository) GetByShipment(ctx context.Context, shipmentID string) (*Estimation, error) {
	e := &Estimation{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, distance_km, base_duration_min, traffic_delay_min,
		weather_delay_min, total_duration_min, eta, confidence, route_hash,
		calculated_at, expires_at FROM estimations
		WHERE shipment_id=$1 ORDER BY calculated_at DESC LIMIT 1`, shipmentID).Scan(
		&e.ID, &e.ShipmentID, &e.DistanceKm, &e.BaseDurationMin, &e.TrafficDelayMin,
		&e.WeatherDelayMin, &e.TotalDurationMin, &e.ETA, &e.Confidence, &e.RouteHash,
		&e.CalculatedAt, &e.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *PostgresRepository) GetLatestByShipment(ctx context.Context, shipmentID string) (*Estimation, error) {
	return r.GetByShipment(ctx, shipmentID)
}
