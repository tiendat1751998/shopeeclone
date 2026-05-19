package fulfillment

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

func (r *PostgresRepository) Create(ctx context.Context, f *Fulfillment) error {
	query := `INSERT INTO fulfillments (
		id, shipment_id, order_id, warehouse_id, status,
		packed_at, shipped_at, completed_at, notes, replay_id, created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, query,
		f.ID, f.ShipmentID, f.OrderID, f.WarehouseID, f.Status,
		f.PackedAt, f.ShippedAt, f.CompletedAt, f.Notes, f.ReplayID,
		f.CreatedAt, f.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Fulfillment, error) {
	f := &Fulfillment{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, order_id, warehouse_id, status,
		packed_at, shipped_at, completed_at, notes, replay_id, created_at, updated_at
		FROM fulfillments WHERE id=$1`, id).Scan(
		&f.ID, &f.ShipmentID, &f.OrderID, &f.WarehouseID, &f.Status,
		&f.PackedAt, &f.ShippedAt, &f.CompletedAt, &f.Notes, &f.ReplayID,
		&f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *PostgresRepository) GetByShipment(ctx context.Context, shipmentID string) (*Fulfillment, error) {
	f := &Fulfillment{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, order_id, warehouse_id, status,
		packed_at, shipped_at, completed_at, notes, replay_id, created_at, updated_at
		FROM fulfillments WHERE shipment_id=$1 ORDER BY created_at DESC LIMIT 1`, shipmentID).Scan(
		&f.ID, &f.ShipmentID, &f.OrderID, &f.WarehouseID, &f.Status,
		&f.PackedAt, &f.ShippedAt, &f.CompletedAt, &f.Notes, &f.ReplayID,
		&f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *PostgresRepository) Update(ctx context.Context, f *Fulfillment) error {
	f.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, `UPDATE fulfillments SET status=$1, packed_at=$2, shipped_at=$3,
		completed_at=$4, notes=$5, updated_at=$6 WHERE id=$7`,
		f.Status, f.PackedAt, f.ShippedAt, f.CompletedAt, f.Notes, f.UpdatedAt, f.ID)
	return err
}

func (r *PostgresRepository) MarkPacked(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, "UPDATE fulfillments SET status=$1, packed_at=$2, updated_at=$3 WHERE id=$4",
		FulfillmentPacked, now, now, id)
	return err
}

func (r *PostgresRepository) MarkShipped(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, "UPDATE fulfillments SET status=$1, shipped_at=$2, updated_at=$3 WHERE id=$4",
		FulfillmentShipped, now, now, id)
	return err
}

func (r *PostgresRepository) MarkCompleted(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, "UPDATE fulfillments SET status=$1, completed_at=$2, updated_at=$3 WHERE id=$4",
		FulfillmentCompleted, now, now, id)
	return err
}
