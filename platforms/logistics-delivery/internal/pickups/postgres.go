package pickups

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, p *Pickup) error {
	query := `INSERT INTO pickups (
		id, shipment_id, fulfillment_id, courier_id, status,
		address, latitude, longitude, scheduled_at, picked_up_at,
		notes, created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, query,
		p.ID, p.ShipmentID, p.FulfillmentID, p.CourierID, p.Status,
		p.Address, p.Latitude, p.Longitude, p.ScheduledAt, p.PickedUpAt,
		p.Notes, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Pickup, error) {
	p := &Pickup{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, fulfillment_id, courier_id, status,
		address, latitude, longitude, scheduled_at, picked_up_at,
		notes, created_at, updated_at FROM pickups WHERE id=$1`, id).Scan(
		&p.ID, &p.ShipmentID, &p.FulfillmentID, &p.CourierID, &p.Status,
		&p.Address, &p.Latitude, &p.Longitude, &p.ScheduledAt, &p.PickedUpAt,
		&p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostgresRepository) GetByShipment(ctx context.Context, shipmentID string) (*Pickup, error) {
	p := &Pickup{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, fulfillment_id, courier_id, status,
		address, latitude, longitude, scheduled_at, picked_up_at,
		notes, created_at, updated_at
		FROM pickups WHERE shipment_id=$1 ORDER BY created_at DESC LIMIT 1`, shipmentID).Scan(
		&p.ID, &p.ShipmentID, &p.FulfillmentID, &p.CourierID, &p.Status,
		&p.Address, &p.Latitude, &p.Longitude, &p.ScheduledAt, &p.PickedUpAt,
		&p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostgresRepository) Update(ctx context.Context, p *Pickup) error {
	p.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, `UPDATE pickups SET courier_id=$1, status=$2, address=$3,
		latitude=$4, longitude=$5, scheduled_at=$6, picked_up_at=$7, notes=$8, updated_at=$9
		WHERE id=$10`,
		p.CourierID, p.Status, p.Address, p.Latitude, p.Longitude,
		p.ScheduledAt, p.PickedUpAt, p.Notes, p.UpdatedAt, p.ID)
	return err
}

func (r *PostgresRepository) MarkCompleted(ctx context.Context, id string, pickedUpAt time.Time) error {
	_, err := r.pool.Exec(ctx, "UPDATE pickups SET status=$1, picked_up_at=$2, updated_at=NOW() WHERE id=$3",
		PickupCompleted, pickedUpAt, id)
	return err
}

func (r *PostgresRepository) MarkFailed(ctx context.Context, id string, reason string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, "UPDATE pickups SET status=$1, notes=$2, updated_at=$3 WHERE id=$4",
		PickupFailed, reason, now, id)
	return err
}
