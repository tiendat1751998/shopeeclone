package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopee-clone/shopee/services/shipment/internal/domain"
)

type ShipmentRepository struct {
	db *sqlx.DB
}

func NewShipmentRepository(db *sqlx.DB) *ShipmentRepository { return &ShipmentRepository{db: db} }

func (r *ShipmentRepository) Create(ctx context.Context, s *domain.Shipment) error {
	query := `INSERT INTO shipments (id, order_id, user_id, carrier_id, tracking_number, status, origin_address, destination_address, weight, dimensions, label_url, cost, currency, idempotency_key, metadata, version, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.OrderID, s.UserID, s.CarrierID, s.TrackingNumber, s.Status, s.OriginAddress, s.DestAddress, s.Weight, s.Dimensions, s.LabelURL, s.Cost, s.Currency, s.IdempotencyKey, s.Metadata, s.Version, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *ShipmentRepository) FindByID(ctx context.Context, id string) (*domain.Shipment, error) {
	var s domain.Shipment
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM shipments WHERE id = ? AND deleted_at IS NULL", id); err != nil {
		if err == sql.ErrNoRows { return nil, domain.ErrShipmentNotFound }
		return nil, err
	}
	return &s, nil
}

func (r *ShipmentRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Shipment, error) {
	var s domain.Shipment
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM shipments WHERE order_id = ? AND deleted_at IS NULL", orderID); err != nil {
		if err == sql.ErrNoRows { return nil, domain.ErrShipmentNotFound }
		return nil, err
	}
	return &s, nil
}

func (r *ShipmentRepository) UpdateStatus(ctx context.Context, id string, status domain.ShipmentStatus, version int) error {
	result, err := r.db.ExecContext(ctx, "UPDATE shipments SET status = ?, version = version + 1, updated_at = ? WHERE id = ? AND version = ?", status, time.Now().UTC(), id, version)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return domain.ErrConcurrentModification }
	return nil
}

func (r *ShipmentRepository) Update(ctx context.Context, s *domain.Shipment) error {
	query := `UPDATE shipments SET status = ?, tracking_number = ?, label_url = ?, cost = ?, metadata = ?, version = version + 1, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, s.Status, s.TrackingNumber, s.LabelURL, s.Cost, s.Metadata, time.Now().UTC(), s.ID)
	return err
}

func (r *ShipmentRepository) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Shipment, error) {
	var s domain.Shipment
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM shipments WHERE idempotency_key = ?", key); err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &s, nil
}

func (r *ShipmentRepository) SaveTrackingEvent(ctx context.Context, event *domain.TrackingEvent) error {
	query := `INSERT INTO tracking_events (id, shipment_id, status, location, description, timestamp, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, event.ID, event.ShipmentID, event.Status, event.Location, event.Description, event.Timestamp, event.CreatedAt)
	return err
}

func (r *ShipmentRepository) GetTrackingHistory(ctx context.Context, shipmentID string) ([]*domain.TrackingEvent, error) {
	var events []*domain.TrackingEvent
	err := r.db.SelectContext(ctx, &events, "SELECT * FROM tracking_events WHERE shipment_id = ? ORDER BY timestamp ASC LIMIT 200", shipmentID)
	return events, err
}

func (r *ShipmentRepository) SaveOutboxEvent(ctx context.Context, event *domain.OutboxEvent) error {
	query := `INSERT INTO outbox_events (event_id, aggregate_type, aggregate_id, event_type, payload, created_at, processed) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, event.ID, event.AggregateType, event.AggregateID, event.EventType, event.Payload, event.CreatedAt, event.Processed)
	return err
}

func (r *ShipmentRepository) GetUnprocessedOutboxEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	var events []*domain.OutboxEvent
	err := r.db.SelectContext(ctx, &events, "SELECT * FROM outbox_events WHERE processed = FALSE ORDER BY created_at ASC LIMIT ?", limit)
	return events, err
}

func (r *ShipmentRepository) MarkOutboxEventProcessed(ctx context.Context, eventID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE outbox_events SET processed = TRUE WHERE event_id = ?", eventID)
	return err
}

func (r *ShipmentRepository) SaveIdempotencyKey(ctx context.Context, record *domain.IdempotencyRecord) error {
	query := `INSERT INTO idempotency_keys (` + "`key`" + `, shipment_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, record.Key, record.ShipmentID, record.ExpiresAt, record.CreatedAt)
	return err
}

func (r *ShipmentRepository) GetIdempotencyKey(ctx context.Context, key string) (*domain.IdempotencyRecord, error) {
	var record domain.IdempotencyRecord
	if err := r.db.GetContext(ctx, &record, "SELECT * FROM idempotency_keys WHERE ` + "`key`" + ` = ?", key); err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &record, nil
}

func (r *ShipmentRepository) IsWebhookProcessed(ctx context.Context, key string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM webhook_events WHERE idempotency_key = ?", key)
	return count > 0, err
}

func (r *ShipmentRepository) SaveWebhookEvent(ctx context.Context, provider, eventType string, payload []byte, signature, idempotencyKey string) error {
	query := `INSERT INTO webhook_events (id, provider, event_type, payload, signature, idempotency_key, created_at) VALUES (?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query, idempotencyKey[:36], provider, eventType, payload, signature, idempotencyKey)
	return err
}
