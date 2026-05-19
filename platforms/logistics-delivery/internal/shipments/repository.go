package shipments

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, s *Shipment) error
	GetByID(ctx context.Context, id string) (*Shipment, error)
	Update(ctx context.Context, s *Shipment) error
	List(ctx context.Context, filter ShipmentFilter) ([]*Shipment, int64, error)
	Delete(ctx context.Context, id string) error
	TransitionStatus(ctx context.Context, txnID string, from, to ShipmentStatus, replayID string) error
	GetByOrderID(ctx context.Context, orderID string) ([]*Shipment, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, s *Shipment) error {
	query := `INSERT INTO shipments (
		id, order_id, customer_id, warehouse_id, courier_id, status,
		origin_street, origin_city, origin_state, origin_country, origin_zip,
		origin_lat, origin_lng,
		dest_street, dest_city, dest_state, dest_country, dest_zip,
		dest_lat, dest_lng,
		total_weight, total_volume, estimated_distance, estimated_eta,
		courier_notes, replay_id, version, created_at, updated_at
	) VALUES (
		$1,$2,$3,$4,$5,$6,
		$7,$8,$9,$10,$11,$12,$13,
		$14,$15,$16,$17,$18,$19,$20,
		$21,$22,$23,$24,$25,$26,1,$27,$28
	)`
	_, err := r.pool.Exec(ctx, query,
		s.ID, s.OrderID, s.CustomerID, s.WarehouseID, s.CourierID, s.Status,
		s.OriginAddress.Street, s.OriginAddress.City, s.OriginAddress.State, s.OriginAddress.Country, s.OriginAddress.ZipCode,
		s.OriginAddress.Latitude, s.OriginAddress.Longitude,
		s.DestinationAddress.Street, s.DestinationAddress.City, s.DestinationAddress.State, s.DestinationAddress.Country, s.DestinationAddress.ZipCode,
		s.DestinationAddress.Latitude, s.DestinationAddress.Longitude,
		s.TotalWeight, s.TotalVolume, s.EstimatedDistance, s.EstimatedETA,
		s.CourierNotes, s.ReplayID, s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Shipment, error) {
	query := `SELECT
		id, order_id, customer_id, warehouse_id, courier_id, status,
		origin_street, origin_city, origin_state, origin_country, origin_zip,
		origin_lat, origin_lng,
		dest_street, dest_city, dest_state, dest_country, dest_zip,
		dest_lat, dest_lng,
		total_weight, total_volume, estimated_distance, estimated_eta,
		actual_delivered_at, courier_notes, replay_id, version, created_at, updated_at
	FROM shipments WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	s := &Shipment{}
	var estETA, actDelivered *time.Time
	err := row.Scan(
		&s.ID, &s.OrderID, &s.CustomerID, &s.WarehouseID, &s.CourierID, &s.Status,
		&s.OriginAddress.Street, &s.OriginAddress.City, &s.OriginAddress.State, &s.OriginAddress.Country, &s.OriginAddress.ZipCode,
		&s.OriginAddress.Latitude, &s.OriginAddress.Longitude,
		&s.DestinationAddress.Street, &s.DestinationAddress.City, &s.DestinationAddress.State, &s.DestinationAddress.Country, &s.DestinationAddress.ZipCode,
		&s.DestinationAddress.Latitude, &s.DestinationAddress.Longitude,
		&s.TotalWeight, &s.TotalVolume, &s.EstimatedDistance, &estETA,
		&actDelivered, &s.CourierNotes, &s.ReplayID, &s.Version, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	s.EstimatedETA = estETA
	s.ActualDeliveredAt = actDelivered
	return s, nil
}

func (r *PostgresRepository) Update(ctx context.Context, s *Shipment) error {
	query := `UPDATE shipments SET
		courier_id=$1, status=$2, total_weight=$3, total_volume=$4,
		estimated_distance=$5, estimated_eta=$6, courier_notes=$7,
		replay_id=$8, version=version+1, updated_at=NOW()
	WHERE id=$9 AND version=$10`
	res, err := r.pool.Exec(ctx, query,
		s.CourierID, s.Status, s.TotalWeight, s.TotalVolume,
		s.EstimatedDistance, s.EstimatedETA, s.CourierNotes,
		s.ReplayID, s.ID, s.Version,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrShipmentNotFound
	}
	return nil
}

func (r *PostgresRepository) List(ctx context.Context, filter ShipmentFilter) ([]*Shipment, int64, error) {
	where := " WHERE 1=1"
	args := []any{}
	argIdx := 1
	if filter.Status != "" {
		where += " AND status=$" + itoa(argIdx); argIdx++
		args = append(args, filter.Status)
	}
	if filter.CourierID != "" {
		where += " AND courier_id=$" + itoa(argIdx); argIdx++
		args = append(args, filter.CourierID)
	}
	if filter.CustomerID != "" {
		where += " AND customer_id=$" + itoa(argIdx); argIdx++
		args = append(args, filter.CustomerID)
	}
	if filter.OrderID != "" {
		where += " AND order_id=$" + itoa(argIdx); argIdx++
		args = append(args, filter.OrderID)
	}
	if filter.WarehouseID != "" {
		where += " AND warehouse_id=$" + itoa(argIdx); argIdx++
		args = append(args, filter.WarehouseID)
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	countQuery := "SELECT COUNT(*) FROM shipments" + where
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	dataQuery := "SELECT id, order_id, customer_id, warehouse_id, courier_id, status, created_at, updated_at FROM shipments" + where +
		" ORDER BY created_at DESC LIMIT $" + itoa(argIdx) + " OFFSET $" + itoa(argIdx+1)
	args = append(args, filter.Limit, filter.Offset)
	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []*Shipment
	for rows.Next() {
		s := &Shipment{}
		if err := rows.Scan(&s.ID, &s.OrderID, &s.CustomerID, &s.WarehouseID, &s.CourierID, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, s)
	}
	return result, total, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM shipments WHERE id=$1", id)
	return err
}

func (r *PostgresRepository) TransitionStatus(ctx context.Context, txnID string, from, to ShipmentStatus, replayID string) error {
	query := `UPDATE shipments SET status=$1, version=version+1, updated_at=NOW()
		WHERE id=$2 AND status=$3 AND (replay_id IS NULL OR replay_id!=$4)`
	res, err := r.pool.Exec(ctx, query, to, txnID, from, replayID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrInvalidStatusTransition
	}
	return nil
}

func (r *PostgresRepository) GetByOrderID(ctx context.Context, orderID string) ([]*Shipment, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, order_id, status, created_at, updated_at FROM shipments WHERE order_id=$1 ORDER BY created_at DESC", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*Shipment
	for rows.Next() {
		s := &Shipment{}
		if err := rows.Scan(&s.ID, &s.OrderID, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0' + n/10)) + string(rune('0' + n%10))
}
