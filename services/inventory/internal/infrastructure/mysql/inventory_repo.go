package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopee-clone/shopee/services/inventory/internal/domain"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) *InventoryRepository { return &InventoryRepository{db: db} }

func (r *InventoryRepository) GetStock(ctx context.Context, skuID, warehouseID string) (*domain.Stock, error) {
	var s domain.Stock
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM stock WHERE sku_id = ? AND warehouse_id = ?", skuID, warehouseID); err != nil {
		if err == sql.ErrNoRows { return nil, domain.ErrStockNotFound }
		return nil, err
	}
	return &s, nil
}

func (r *InventoryRepository) UpdateStock(ctx context.Context, stock *domain.Stock) error {
	query := `UPDATE stock SET quantity = ?, reserved_qty = ?, available_qty = ?, status = ?, version = version + 1, updated_at = ? WHERE id = ? AND version = ?`
	result, err := r.db.ExecContext(ctx, query, stock.Quantity, stock.ReservedQty, stock.AvailableQty, stock.Status, time.Now().UTC(), stock.ID, stock.Version)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return domain.ErrConcurrentModification }
	return nil
}

func (r *InventoryRepository) CreateStock(ctx context.Context, stock *domain.Stock) error {
	query := `INSERT INTO stock (id, product_id, sku_id, warehouse_id, quantity, reserved_qty, available_qty, status, reorder_level, version, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, stock.ID, stock.ProductID, stock.SkuID, stock.WarehouseID, stock.Quantity, stock.ReservedQty, stock.AvailableQty, stock.Status, stock.ReorderLevel, stock.Version, stock.CreatedAt, stock.UpdatedAt)
	return err
}

func (r *InventoryRepository) SaveReservation(ctx context.Context, res *domain.Reservation) error {
	query := `INSERT INTO reservations (id, order_id, user_id, product_id, sku_id, warehouse_id, quantity, status, expires_at, idempotency_key, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, res.ID, res.OrderID, res.UserID, res.ProductID, res.SkuID, res.WarehouseID, res.Quantity, res.Status, res.ExpiresAt, res.IdempotencyKey, res.CreatedAt, res.UpdatedAt)
	return err
}

func (r *InventoryRepository) GetReservation(ctx context.Context, id string) (*domain.Reservation, error) {
	var res domain.Reservation
	if err := r.db.GetContext(ctx, &res, "SELECT * FROM reservations WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows { return nil, domain.ErrReservationNotFound }
		return nil, err
	}
	return &res, nil
}

func (r *InventoryRepository) UpdateReservationStatus(ctx context.Context, id string, status domain.ReservationStatus) error {
	_, err := r.db.ExecContext(ctx, "UPDATE reservations SET status = ?, updated_at = ? WHERE id = ?", status, time.Now().UTC(), id)
	return err
}

func (r *InventoryRepository) GetExpiredReservations(ctx context.Context, limit int) ([]*domain.Reservation, error) {
	var res []*domain.Reservation
	err := r.db.SelectContext(ctx, &res, "SELECT * FROM reservations WHERE status = 'active' AND expires_at < NOW() LIMIT ?", limit)
	return res, err
}

func (r *InventoryRepository) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Reservation, error) {
	var res domain.Reservation
	if err := r.db.GetContext(ctx, &res, "SELECT * FROM reservations WHERE idempotency_key = ?", key); err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &res, nil
}

func (r *InventoryRepository) SaveOutboxEvent(ctx context.Context, event *domain.OutboxEvent) error {
	query := `INSERT INTO outbox_events (event_id, aggregate_type, aggregate_id, event_type, payload, created_at, processed) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, event.ID, event.AggregateType, event.AggregateID, event.EventType, event.Payload, event.CreatedAt, event.Processed)
	return err
}

func (r *InventoryRepository) GetUnprocessedOutboxEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	var events []*domain.OutboxEvent
	err := r.db.SelectContext(ctx, &events, "SELECT * FROM outbox_events WHERE processed = FALSE ORDER BY created_at ASC LIMIT ?", limit)
	return events, err
}

func (r *InventoryRepository) MarkOutboxEventProcessed(ctx context.Context, eventID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE outbox_events SET processed = TRUE WHERE event_id = ?", eventID)
	return err
}

func (r *InventoryRepository) SaveIdempotencyKey(ctx context.Context, record *domain.IdempotencyRecord) error {
	query := `INSERT INTO idempotency_keys (` + "`key`" + `, reservation_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, record.Key, record.ReservationID, record.ExpiresAt, record.CreatedAt)
	return err
}

func (r *InventoryRepository) GetIdempotencyKey(ctx context.Context, key string) (*domain.IdempotencyRecord, error) {
	var record domain.IdempotencyRecord
	if err := r.db.GetContext(ctx, &record, "SELECT * FROM idempotency_keys WHERE ` + "`key`" + ` = ?", key); err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &record, nil
}
