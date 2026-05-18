package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/shopee-clone/shopee/services/inventory/internal/domain"
)

type StockRepository struct {
	db *sqlx.DB
}

func NewStockRepository(db *sqlx.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) FindBySKUAndWarehouse(ctx context.Context, sku, warehouseID string) (*domain.Stock, error) {
	var stock domain.Stock
	query := `SELECT id, sku, warehouse_id, quantity, reserved, available, 
		flash_sale_stock, flash_sale_reserved, version, created_at, updated_at 
		FROM stocks WHERE sku = ? AND warehouse_id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &stock, query, sku, warehouseID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find stock: %w", err)
	}
	return &stock, nil
}

func (r *StockRepository) FindBySKU(ctx context.Context, sku string) ([]*domain.Stock, error) {
	var stocks []*domain.Stock
	query := `SELECT id, sku, warehouse_id, quantity, reserved, available, 
		flash_sale_stock, flash_sale_reserved, version, created_at, updated_at 
		FROM stocks WHERE sku = ? AND deleted_at IS NULL ORDER BY warehouse_id`
	err := r.db.SelectContext(ctx, &stocks, query, sku)
	if err != nil {
		return nil, fmt.Errorf("find stocks by sku: %w", err)
	}
	return stocks, nil
}

func (r *StockRepository) FindByWarehouse(ctx context.Context, warehouseID string, offset, limit int) ([]*domain.Stock, error) {
	var stocks []*domain.Stock
	query := `SELECT id, sku, warehouse_id, quantity, reserved, available, 
		flash_sale_stock, flash_sale_reserved, version, created_at, updated_at 
		FROM stocks WHERE warehouse_id = ? AND deleted_at IS NULL 
		ORDER BY sku LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &stocks, query, warehouseID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("find stocks by warehouse: %w", err)
	}
	return stocks, nil
}

func (r *StockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	query := `INSERT INTO stocks (id, sku, warehouse_id, quantity, reserved, available, 
		flash_sale_stock, flash_sale_reserved, version, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		stock.ID, stock.SKU, stock.WarehouseID, stock.Quantity, stock.Reserved,
		stock.Available, stock.FlashSaleStock, stock.FlashSaleReserved,
		stock.Version, stock.CreatedAt, stock.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create stock: %w", err)
	}
	return nil
}

func (r *StockRepository) Update(ctx context.Context, stock *domain.Stock) error {
	query := `UPDATE stocks SET quantity = ?, reserved = ?, available = ?, 
		flash_sale_stock = ?, flash_sale_reserved = ?, version = version + 1, updated_at = ? 
		WHERE id = ? AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query,
		stock.Quantity, stock.Reserved, stock.Available,
		stock.FlashSaleStock, stock.FlashSaleReserved, stock.UpdatedAt, stock.ID)
	if err != nil {
		return fmt.Errorf("update stock: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("stock not found: %s", stock.ID)
	}
	stock.Version++
	return nil
}

func (r *StockRepository) UpdateWithVersion(ctx context.Context, stock *domain.Stock) error {
	query := `UPDATE stocks SET quantity = ?, reserved = ?, available = ?, 
		flash_sale_stock = ?, flash_sale_reserved = ?, version = version + 1, updated_at = ? 
		WHERE id = ? AND version = ? AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query,
		stock.Quantity, stock.Reserved, stock.Available,
		stock.FlashSaleStock, stock.FlashSaleReserved, stock.UpdatedAt,
		stock.ID, stock.Version)
	if err != nil {
		return fmt.Errorf("update stock with version: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("optimistic lock conflict for stock %s version %d", stock.ID, stock.Version)
	}
	stock.Version++
	return nil
}

func (r *StockRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE stocks SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *StockRepository) List(ctx context.Context, filter domain.StockFilter) ([]*domain.Stock, int64, error) {
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "deleted_at IS NULL")

	if filter.SKU != "" {
		conditions = append(conditions, "sku = ?")
		args = append(args, filter.SKU)
	}
	if filter.WarehouseID != "" {
		conditions = append(conditions, "warehouse_id = ?")
		args = append(args, filter.WarehouseID)
	}
	if filter.MinAvail != nil {
		conditions = append(conditions, "available >= ?")
		args = append(args, *filter.MinAvail)
	}
	if filter.MaxAvail != nil {
		conditions = append(conditions, "available <= ?")
		args = append(args, *filter.MaxAvail)
	}

	whereClause := strings.Join(conditions, " AND ")

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM stocks WHERE %s", whereClause)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count stocks: %w", err)
	}

	if filter.Limit == 0 {
		filter.Limit = 20
	}

	selectQuery := fmt.Sprintf(`SELECT id, sku, warehouse_id, quantity, reserved, available, 
		flash_sale_stock, flash_sale_reserved, version, created_at, updated_at 
		FROM stocks WHERE %s ORDER BY updated_at DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, filter.Limit, filter.Offset)

	var stocks []*domain.Stock
	if err := r.db.SelectContext(ctx, &stocks, selectQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("list stocks: %w", err)
	}

	return stocks, total, nil
}
