package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shopee-clone/shopee/services/inventory/internal/domain"
)

type WarehouseRepository struct {
	db *sqlx.DB
}

func NewWarehouseRepository(db *sqlx.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) FindByID(ctx context.Context, id string) (*domain.Warehouse, error) {
	var w domain.Warehouse
	query := `SELECT id, name, code, address, city, region, priority, is_active, created_at, updated_at 
		FROM warehouses WHERE id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &w, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find warehouse: %w", err)
	}
	return &w, nil
}

func (r *WarehouseRepository) FindByCode(ctx context.Context, code string) (*domain.Warehouse, error) {
	var w domain.Warehouse
	query := `SELECT id, name, code, address, city, region, priority, is_active, created_at, updated_at 
		FROM warehouses WHERE code = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &w, query, code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find warehouse by code: %w", err)
	}
	return &w, nil
}

func (r *WarehouseRepository) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `INSERT INTO warehouses (id, name, code, address, city, region, priority, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		warehouse.ID, warehouse.Name, warehouse.Code, warehouse.Address,
		warehouse.City, warehouse.Region, warehouse.Priority, warehouse.IsActive,
		warehouse.CreatedAt, warehouse.UpdatedAt)
	return err
}

func (r *WarehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `UPDATE warehouses SET name = ?, address = ?, city = ?, region = ?, 
		priority = ?, is_active = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query,
		warehouse.Name, warehouse.Address, warehouse.City, warehouse.Region,
		warehouse.Priority, warehouse.IsActive, warehouse.UpdatedAt, warehouse.ID)
	return err
}

func (r *WarehouseRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE warehouses SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	return err
}

func (r *WarehouseRepository) ListActive(ctx context.Context) ([]*domain.Warehouse, error) {
	var warehouses []*domain.Warehouse
	query := `SELECT id, name, code, address, city, region, priority, is_active, created_at, updated_at 
		FROM warehouses WHERE is_active = true AND deleted_at IS NULL ORDER BY priority ASC`
	err := r.db.SelectContext(ctx, &warehouses, query)
	return warehouses, err
}

func (r *WarehouseRepository) List(ctx context.Context, offset, limit int) ([]*domain.Warehouse, int64, error) {
	var total int64
	if err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM warehouses WHERE deleted_at IS NULL"); err != nil {
		return nil, 0, err
	}

	var warehouses []*domain.Warehouse
	query := `SELECT id, name, code, address, city, region, priority, is_active, created_at, updated_at 
		FROM warehouses WHERE deleted_at IS NULL ORDER BY priority ASC LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &warehouses, query, limit, offset)
	return warehouses, total, err
}

type StockMovementRepository struct {
	db *sqlx.DB
}

func NewStockMovementRepository(db *sqlx.DB) *StockMovementRepository {
	return &StockMovementRepository{db: db}
}

func (r *StockMovementRepository) Create(ctx context.Context, movement *domain.StockMovement) error {
	query := `INSERT INTO stock_movements (id, sku, warehouse_id, type, quantity, before_qty, after_qty, reference_id, reason, operator_id, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		movement.ID, movement.SKU, movement.WarehouseID, movement.Type,
		movement.Quantity, movement.BeforeQty, movement.AfterQty,
		movement.ReferenceID, movement.Reason, movement.OperatorID, movement.CreatedAt)
	return err
}

func (r *StockMovementRepository) FindBySKU(ctx context.Context, sku string, offset, limit int) ([]*domain.StockMovement, int64, error) {
	var total int64
	r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM stock_movements WHERE sku = ?", sku)

	var movements []*domain.StockMovement
	query := `SELECT id, sku, warehouse_id, type, quantity, before_qty, after_qty, reference_id, reason, operator_id, created_at 
		FROM stock_movements WHERE sku = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &movements, query, sku, limit, offset)
	return movements, total, err
}

func (r *StockMovementRepository) FindByWarehouse(ctx context.Context, warehouseID string, offset, limit int) ([]*domain.StockMovement, int64, error) {
	var total int64
	r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM stock_movements WHERE warehouse_id = ?", warehouseID)

	var movements []*domain.StockMovement
	query := `SELECT id, sku, warehouse_id, type, quantity, before_qty, after_qty, reference_id, reason, operator_id, created_at 
		FROM stock_movements WHERE warehouse_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &movements, query, warehouseID, limit, offset)
	return movements, total, err
}

func (r *StockMovementRepository) FindByReference(ctx context.Context, referenceID string) ([]*domain.StockMovement, error) {
	var movements []*domain.StockMovement
	query := `SELECT id, sku, warehouse_id, type, quantity, before_qty, after_qty, reference_id, reason, operator_id, created_at 
		FROM stock_movements WHERE reference_id = ? ORDER BY created_at DESC LIMIT 100`
	err := r.db.SelectContext(ctx, &movements, query, referenceID)
	return movements, err
}

type FlashSaleInventoryRepository struct {
	db *sqlx.DB
}

func NewFlashSaleInventoryRepository(db *sqlx.DB) *FlashSaleInventoryRepository {
	return &FlashSaleInventoryRepository{db: db}
}

func (r *FlashSaleInventoryRepository) FindByID(ctx context.Context, id string) (*domain.FlashSaleInventory, error) {
	var fs domain.FlashSaleInventory
	query := `SELECT id, flash_sale_id, sku, warehouse_id, total_stock, reserved_stock, sold_stock, 
		max_per_user, start_time, end_time, is_active, created_at, updated_at 
		FROM flash_sale_inventory WHERE id = ?`
	err := r.db.GetContext(ctx, &fs, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fs, nil
}

func (r *FlashSaleInventoryRepository) FindByFlashSaleAndSKU(ctx context.Context, flashSaleID, sku string) (*domain.FlashSaleInventory, error) {
	var fs domain.FlashSaleInventory
	query := `SELECT id, flash_sale_id, sku, warehouse_id, total_stock, reserved_stock, sold_stock, 
		max_per_user, start_time, end_time, is_active, created_at, updated_at 
		FROM flash_sale_inventory WHERE flash_sale_id = ? AND sku = ?`
	err := r.db.GetContext(ctx, &fs, query, flashSaleID, sku)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fs, nil
}

func (r *FlashSaleInventoryRepository) FindByFlashSale(ctx context.Context, flashSaleID string) ([]*domain.FlashSaleInventory, error) {
	var items []*domain.FlashSaleInventory
	query := `SELECT id, flash_sale_id, sku, warehouse_id, total_stock, reserved_stock, sold_stock, 
		max_per_user, start_time, end_time, is_active, created_at, updated_at 
		FROM flash_sale_inventory WHERE flash_sale_id = ?`
	err := r.db.SelectContext(ctx, &items, query, flashSaleID)
	return items, err
}

func (r *FlashSaleInventoryRepository) Create(ctx context.Context, fs *domain.FlashSaleInventory) error {
	query := `INSERT INTO flash_sale_inventory (id, flash_sale_id, sku, warehouse_id, total_stock, 
		reserved_stock, sold_stock, max_per_user, start_time, end_time, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		fs.ID, fs.FlashSaleID, fs.SKU, fs.WarehouseID, fs.TotalStock,
		fs.ReservedStock, fs.SoldStock, fs.MaxPerUser, fs.StartTime,
		fs.EndTime, fs.IsActive, fs.CreatedAt, fs.UpdatedAt)
	return err
}

func (r *FlashSaleInventoryRepository) Update(ctx context.Context, fs *domain.FlashSaleInventory) error {
	query := `UPDATE flash_sale_inventory SET total_stock = ?, reserved_stock = ?, sold_stock = ?, 
		max_per_user = ?, is_active = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		fs.TotalStock, fs.ReservedStock, fs.SoldStock, fs.MaxPerUser,
		fs.IsActive, fs.UpdatedAt, fs.ID)
	return err
}

func (r *FlashSaleInventoryRepository) UpdateStock(ctx context.Context, id string, reservedDelta, soldDelta int64) error {
	query := `UPDATE flash_sale_inventory SET reserved_stock = reserved_stock + ?, 
		sold_stock = sold_stock + ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, reservedDelta, soldDelta, id)
	return err
}

func (r *FlashSaleInventoryRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM flash_sale_inventory WHERE id = ?", id)
	return err
}

func (r *FlashSaleInventoryRepository) ListActive(ctx context.Context) ([]*domain.FlashSaleInventory, error) {
	var items []*domain.FlashSaleInventory
	query := `SELECT id, flash_sale_id, sku, warehouse_id, total_stock, reserved_stock, sold_stock, 
		max_per_user, start_time, end_time, is_active, created_at, updated_at 
		FROM flash_sale_inventory WHERE is_active = true AND end_time > NOW() 
		ORDER BY start_time ASC`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}
