package dispatch

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

func (r *PostgresRepository) Create(ctx context.Context, d *Dispatch) error {
	query := `INSERT INTO dispatches (
		id, shipment_id, courier_id, zone_id, status,
		pickup_time, dispatch_time, completed_at, notes, created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, query,
		d.ID, d.ShipmentID, d.CourierID, d.ZoneID, d.Status,
		d.PickupTime, d.DispatchTime, d.CompletedAt, d.Notes,
		d.CreatedAt, d.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Dispatch, error) {
	d := &Dispatch{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, courier_id, zone_id, status,
		pickup_time, dispatch_time, completed_at, notes, created_at, updated_at
		FROM dispatches WHERE id=$1`, id).Scan(
		&d.ID, &d.ShipmentID, &d.CourierID, &d.ZoneID, &d.Status,
		&d.PickupTime, &d.DispatchTime, &d.CompletedAt, &d.Notes,
		&d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *PostgresRepository) GetByShipment(ctx context.Context, shipmentID string) (*Dispatch, error) {
	d := &Dispatch{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, courier_id, zone_id, status,
		pickup_time, dispatch_time, completed_at, notes, created_at, updated_at
		FROM dispatches WHERE shipment_id=$1 ORDER BY created_at DESC LIMIT 1`, shipmentID).Scan(
		&d.ID, &d.ShipmentID, &d.CourierID, &d.ZoneID, &d.Status,
		&d.PickupTime, &d.DispatchTime, &d.CompletedAt, &d.Notes,
		&d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *PostgresRepository) Update(ctx context.Context, d *Dispatch) error {
	d.UpdatedAt = time.Now().UTC()
	query := `UPDATE dispatches SET
		courier_id=$1, zone_id=$2, status=$3, pickup_time=$4,
		dispatch_time=$5, completed_at=$6, notes=$7, updated_at=$8
		WHERE id=$9`
	_, err := r.pool.Exec(ctx, query,
		d.CourierID, d.ZoneID, d.Status, d.PickupTime,
		d.DispatchTime, d.CompletedAt, d.Notes, d.UpdatedAt, d.ID)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, filter DispatchFilter) ([]*Dispatch, int64, error) {
	where := " WHERE 1=1"
	args := []any{}
	idx := 1
	if filter.Status != "" {
		where += " AND status=$" + itoa(idx); idx++
		args = append(args, filter.Status)
	}
	if filter.CourierID != "" {
		where += " AND courier_id=$" + itoa(idx); idx++
		args = append(args, filter.CourierID)
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM dispatches"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	query := `SELECT id, shipment_id, courier_id, zone_id, status,
		pickup_time, dispatch_time, completed_at, notes, created_at, updated_at
		FROM dispatches` + where + ` ORDER BY created_at DESC LIMIT $` + itoa(idx) + ` OFFSET $` + itoa(idx+1)
	args = append(args, filter.Limit, filter.Offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []*Dispatch
	for rows.Next() {
		d := &Dispatch{}
		if err := rows.Scan(&d.ID, &d.ShipmentID, &d.CourierID, &d.ZoneID, &d.Status,
			&d.PickupTime, &d.DispatchTime, &d.CompletedAt, &d.Notes,
			&d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, d)
	}
	return result, total, nil
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0' + n/10)) + string(rune('0' + n%10))
}
