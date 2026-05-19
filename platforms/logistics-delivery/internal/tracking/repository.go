package tracking

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	AppendEvent(ctx context.Context, e *TrackingEvent) error
	GetTimeline(ctx context.Context, shipmentID string) ([]*TrackingEvent, error)
	GetMilestones(ctx context.Context, shipmentID string) ([]Milestone, error)
	ListEvents(ctx context.Context, filter TrackingFilter) ([]*TrackingEvent, int64, error)
	GetLastEvent(ctx context.Context, shipmentID string) (*TrackingEvent, error)
	DeleteEvents(ctx context.Context, shipmentID string) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) AppendEvent(ctx context.Context, e *TrackingEvent) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	query := `INSERT INTO tracking_events (
		id, shipment_id, event_type, lat, lng, location_name, location_address,
		description, courier_data, replay_id, occurred_at, created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, query,
		e.ID, e.ShipmentID, e.EventType,
		e.Location.Latitude, e.Location.Longitude, e.Location.Name, e.Location.Address,
		e.Description, e.CourierData, e.ReplayID, e.OccurredAt, e.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetTimeline(ctx context.Context, shipmentID string) ([]*TrackingEvent, error) {
	rows, err := r.pool.Query(ctx, `SELECT
		id, shipment_id, event_type, lat, lng, location_name, location_address,
		description, courier_data, replay_id, occurred_at, created_at
		FROM tracking_events WHERE shipment_id=$1 ORDER BY occurred_at ASC LIMIT 200`, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*TrackingEvent
	for rows.Next() {
		e := &TrackingEvent{}
		if err := rows.Scan(
			&e.ID, &e.ShipmentID, &e.EventType,
			&e.Location.Latitude, &e.Location.Longitude, &e.Location.Name, &e.Location.Address,
			&e.Description, &e.CourierData, &e.ReplayID, &e.OccurredAt, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *PostgresRepository) GetMilestones(ctx context.Context, shipmentID string) ([]Milestone, error) {
	rows, err := r.pool.Query(ctx, `SELECT DISTINCT ON (event_type) event_type, occurred_at, description
		FROM tracking_events WHERE shipment_id=$1 ORDER BY event_type, occurred_at DESC`, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var milestones []Milestone
	for rows.Next() {
		var m Milestone
		if err := rows.Scan(&m.EventType, &m.AchievedAt, &m.Description); err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}
	return milestones, nil
}

func (r *PostgresRepository) ListEvents(ctx context.Context, filter TrackingFilter) ([]*TrackingEvent, int64, error) {
	where := " WHERE 1=1"
	args := []any{}
	idx := 1
	if filter.ShipmentID != "" {
		where += " AND shipment_id=$" + itoa(idx); idx++
		args = append(args, filter.ShipmentID)
	}
	if filter.FromDate != nil {
		where += " AND occurred_at>=$" + itoa(idx); idx++
		args = append(args, *filter.FromDate)
	}
	if filter.ToDate != nil {
		where += " AND occurred_at<=$" + itoa(idx); idx++
		args = append(args, *filter.ToDate)
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM tracking_events"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	query := `SELECT id, shipment_id, event_type, lat, lng, location_name, location_address,
		description, courier_data, replay_id, occurred_at, created_at
		FROM tracking_events` + where + ` ORDER BY occurred_at DESC LIMIT $` + itoa(idx) + ` OFFSET $` + itoa(idx+1)
	args = append(args, filter.Limit, filter.Offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []*TrackingEvent
	for rows.Next() {
		e := &TrackingEvent{}
		if err := rows.Scan(
			&e.ID, &e.ShipmentID, &e.EventType,
			&e.Location.Latitude, &e.Location.Longitude, &e.Location.Name, &e.Location.Address,
			&e.Description, &e.CourierData, &e.ReplayID, &e.OccurredAt, &e.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, e)
	}
	return result, total, nil
}

func (r *PostgresRepository) GetLastEvent(ctx context.Context, shipmentID string) (*TrackingEvent, error) {
	e := &TrackingEvent{}
	err := r.pool.QueryRow(ctx, `SELECT
		id, shipment_id, event_type, lat, lng, location_name, location_address,
		description, courier_data, replay_id, occurred_at, created_at
		FROM tracking_events WHERE shipment_id=$1 ORDER BY occurred_at DESC LIMIT 1`, shipmentID).Scan(
		&e.ID, &e.ShipmentID, &e.EventType,
		&e.Location.Latitude, &e.Location.Longitude, &e.Location.Name, &e.Location.Address,
		&e.Description, &e.CourierData, &e.ReplayID, &e.OccurredAt, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *PostgresRepository) DeleteEvents(ctx context.Context, shipmentID string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM tracking_events WHERE shipment_id=$1", shipmentID)
	return err
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0' + n/10)) + string(rune('0' + n%10))
}
