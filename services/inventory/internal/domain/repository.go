package domain

import "context"

// StockRepository defines stock data access
type StockRepository interface {
	FindBySKUAndWarehouse(ctx context.Context, sku, warehouseID string) (*Stock, error)
	FindBySKU(ctx context.Context, sku string) ([]*Stock, error)
	FindByWarehouse(ctx context.Context, warehouseID string, offset, limit int) ([]*Stock, error)
	Create(ctx context.Context, stock *Stock) error
	Update(ctx context.Context, stock *Stock) error
	UpdateWithVersion(ctx context.Context, stock *Stock) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter StockFilter) ([]*Stock, int64, error)
}

// ReservationRepository defines reservation data access
type ReservationRepository interface {
	FindByID(ctx context.Context, id string) (*Reservation, error)
	FindByKey(ctx context.Context, reservationKey string) (*Reservation, error)
	FindByOrderID(ctx context.Context, orderID string) ([]*Reservation, error)
	FindByUserID(ctx context.Context, userID string, status string, offset, limit int) ([]*Reservation, int64, error)
	FindExpired(ctx context.Context, before string, limit int) ([]*Reservation, error)
	Create(ctx context.Context, reservation *Reservation) error
	Update(ctx context.Context, reservation *Reservation) error
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
	CountActiveBySKU(ctx context.Context, sku string) (int64, error)
}

// WarehouseRepository defines warehouse data access
type WarehouseRepository interface {
	FindByID(ctx context.Context, id string) (*Warehouse, error)
	FindByCode(ctx context.Context, code string) (*Warehouse, error)
	Create(ctx context.Context, warehouse *Warehouse) error
	Update(ctx context.Context, warehouse *Warehouse) error
	Delete(ctx context.Context, id string) error
	ListActive(ctx context.Context) ([]*Warehouse, error)
	List(ctx context.Context, offset, limit int) ([]*Warehouse, int64, error)
}

// StockMovementRepository defines movement log data access
type StockMovementRepository interface {
	Create(ctx context.Context, movement *StockMovement) error
	FindBySKU(ctx context.Context, sku string, offset, limit int) ([]*StockMovement, int64, error)
	FindByWarehouse(ctx context.Context, warehouseID string, offset, limit int) ([]*StockMovement, int64, error)
	FindByReference(ctx context.Context, referenceID string) ([]*StockMovement, error)
}

// FlashSaleInventoryRepository defines flash-sale inventory data access
type FlashSaleInventoryRepository interface {
	FindByID(ctx context.Context, id string) (*FlashSaleInventory, error)
	FindByFlashSaleAndSKU(ctx context.Context, flashSaleID, sku string) (*FlashSaleInventory, error)
	FindByFlashSale(ctx context.Context, flashSaleID string) ([]*FlashSaleInventory, error)
	Create(ctx context.Context, fs *FlashSaleInventory) error
	Update(ctx context.Context, fs *FlashSaleInventory) error
	UpdateStock(ctx context.Context, id string, reservedDelta, soldDelta int64) error
	Delete(ctx context.Context, id string) error
	ListActive(ctx context.Context) ([]*FlashSaleInventory, error)
}

// StockFilter for querying stocks
type StockFilter struct {
	SKU         string
	WarehouseID string
	MinAvail    *int64
	MaxAvail    *int64
	Offset      int
	Limit       int
}
