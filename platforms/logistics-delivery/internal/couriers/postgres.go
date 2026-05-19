package couriers

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

func (r *PostgresRepository) Create(ctx context.Context, c *Courier) error {
	query := `INSERT INTO couriers (
		id, name, phone, provider, status, zone_id,
		current_lat, current_lng, max_capacity, current_load, rating,
		is_active, created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.pool.Exec(ctx, query,
		c.ID, c.Name, c.Phone, c.Provider, c.Status, c.ZoneID,
		c.CurrentLat, c.CurrentLng, c.MaxCapacity, c.CurrentLoad, c.Rating,
		c.IsActive, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Courier, error) {
	c := &Courier{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, name, phone, provider, status, zone_id,
		current_lat, current_lng, last_seen_at, max_capacity, current_load, rating,
		is_active, created_at, updated_at FROM couriers WHERE id=$1`, id).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Provider, &c.Status, &c.ZoneID,
		&c.CurrentLat, &c.CurrentLng, &c.LastSeenAt, &c.MaxCapacity, &c.CurrentLoad, &c.Rating,
		&c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *PostgresRepository) Update(ctx context.Context, c *Courier) error {
	c.UpdatedAt = time.Now().UTC()
	query := `UPDATE couriers SET name=$1, phone=$2, status=$3, zone_id=$4,
		current_lat=$5, current_lng=$6, last_seen_at=$7, current_load=$8,
		rating=$9, updated_at=$10 WHERE id=$11`
	_, err := r.pool.Exec(ctx, query,
		c.Name, c.Phone, c.Status, c.ZoneID,
		c.CurrentLat, c.CurrentLng, c.LastSeenAt, c.CurrentLoad,
		c.Rating, c.UpdatedAt, c.ID)
	return err
}

func (r *PostgresRepository) ListAvailable(ctx context.Context, zoneID string) ([]*Courier, error) {
	rows, err := r.pool.Query(ctx, `SELECT
		id, name, phone, provider, status, zone_id,
		current_lat, current_lng, last_seen_at, max_capacity, current_load, rating,
		is_active, created_at, updated_at FROM couriers
		WHERE zone_id=$1 AND status='available' AND is_active=true
		ORDER BY rating DESC`, zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*Courier
	for rows.Next() {
		c := &Courier{}
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.Provider, &c.Status, &c.ZoneID,
			&c.CurrentLat, &c.CurrentLng, &c.LastSeenAt, &c.MaxCapacity, &c.CurrentLoad, &c.Rating,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]*Courier, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM couriers").Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT
		id, name, phone, provider, status, zone_id,
		current_lat, current_lng, last_seen_at, max_capacity, current_load, rating,
		is_active, created_at, updated_at FROM couriers
		ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []*Courier
	for rows.Next() {
		c := &Courier{}
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.Provider, &c.Status, &c.ZoneID,
			&c.CurrentLat, &c.CurrentLng, &c.LastSeenAt, &c.MaxCapacity, &c.CurrentLoad, &c.Rating,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, c)
	}
	return result, total, nil
}

func (r *PostgresRepository) UpdateLocation(ctx context.Context, courierID string, lat, lng float64) error {
	_, err := r.pool.Exec(ctx, `UPDATE couriers SET
		current_lat=$1, current_lng=$2, last_seen_at=NOW(), updated_at=NOW()
		WHERE id=$3`, lat, lng, courierID)
	return err
}
