package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shopee-clone/shopee/services/inventory/internal/domain"
)

type ReservationRepository struct {
	db *sqlx.DB
}

func NewReservationRepository(db *sqlx.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) FindByID(ctx context.Context, id string) (*domain.Reservation, error) {
	var res domain.Reservation
	query := `SELECT id, reservation_key, user_id, order_id, sku, warehouse_id, quantity, 
		status, idempotency_key, expires_at, created_at, updated_at 
		FROM reservations WHERE id = ?`
	err := r.db.GetContext(ctx, &res, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find reservation: %w", err)
	}
	return &res, nil
}

func (r *ReservationRepository) FindByKey(ctx context.Context, reservationKey string) (*domain.Reservation, error) {
	var res domain.Reservation
	query := `SELECT id, reservation_key, user_id, order_id, sku, warehouse_id, quantity, 
		status, idempotency_key, expires_at, created_at, updated_at 
		FROM reservations WHERE reservation_key = ?`
	err := r.db.GetContext(ctx, &res, query, reservationKey)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find reservation by key: %w", err)
	}
	return &res, nil
}

func (r *ReservationRepository) FindByOrderID(ctx context.Context, orderID string) ([]*domain.Reservation, error) {
	var reservations []*domain.Reservation
	query := `SELECT id, reservation_key, user_id, order_id, sku, warehouse_id, quantity, 
		status, idempotency_key, expires_at, created_at, updated_at 
		FROM reservations WHERE order_id = ? ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &reservations, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("find reservations by order: %w", err)
	}
	return reservations, nil
}

func (r *ReservationRepository) FindByUserID(ctx context.Context, userID string, status string, offset, limit int) ([]*domain.Reservation, int64, error) {
	var conditions []string
	var args []interface{}
	conditions = append(conditions, "user_id = ?")
	args = append(args, userID)
	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + joinConditions(conditions, " AND ")
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM reservations %s", whereClause)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count reservations: %w", err)
	}

	selectQuery := fmt.Sprintf(`SELECT id, reservation_key, user_id, order_id, sku, warehouse_id, 
		quantity, status, idempotency_key, expires_at, created_at, updated_at 
		FROM reservations %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, limit, offset)

	var reservations []*domain.Reservation
	if err := r.db.SelectContext(ctx, &reservations, selectQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("find reservations by user: %w", err)
	}

	return reservations, total, nil
}

func (r *ReservationRepository) FindExpired(ctx context.Context, before string, limit int) ([]*domain.Reservation, error) {
	var reservations []*domain.Reservation
	query := `SELECT id, reservation_key, user_id, order_id, sku, warehouse_id, quantity, 
		status, idempotency_key, expires_at, created_at, updated_at 
		FROM reservations WHERE status = 'pending' AND expires_at < ? 
		ORDER BY expires_at ASC LIMIT ?`
	err := r.db.SelectContext(ctx, &reservations, query, before, limit)
	if err != nil {
		return nil, fmt.Errorf("find expired reservations: %w", err)
	}
	return reservations, nil
}

func (r *ReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	query := `INSERT INTO reservations (id, reservation_key, user_id, order_id, sku, warehouse_id, 
		quantity, status, idempotency_key, expires_at, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		reservation.ID, reservation.ReservationKey, reservation.UserID, reservation.OrderID,
		reservation.SKU, reservation.WarehouseID, reservation.Quantity, reservation.Status,
		reservation.IdempotencyKey, reservation.ExpiresAt, reservation.CreatedAt, reservation.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create reservation: %w", err)
	}
	return nil
}

func (r *ReservationRepository) Update(ctx context.Context, reservation *domain.Reservation) error {
	query := `UPDATE reservations SET order_id = ?, status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		reservation.OrderID, reservation.Status, reservation.UpdatedAt, reservation.ID)
	return err
}

func (r *ReservationRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE reservations SET status = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *ReservationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM reservations WHERE id = ?", id)
	return err
}

func (r *ReservationRepository) CountActiveBySKU(ctx context.Context, sku string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM reservations WHERE sku = ? AND status IN ('pending', 'confirmed')`
	err := r.db.GetContext(ctx, &count, query, sku)
	return count, err
}

func joinConditions(conditions []string, sep string) string {
	result := ""
	for i, c := range conditions {
		if i > 0 {
			result += sep
		}
		result += c
	}
	return result
}
