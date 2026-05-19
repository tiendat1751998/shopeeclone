package inventory

import (
	"context"
	"sync"
)

type ReservationRepository interface {
	Create(ctx context.Context, r *InventoryReservation) error
	GetByID(ctx context.Context, id string) (*InventoryReservation, error)
	Update(ctx context.Context, r *InventoryReservation) error
	GetByOrderID(ctx context.Context, orderID string) ([]*InventoryReservation, error)
}

type StockRepository interface {
	Get(ctx context.Context, productID, warehouseID string) (*Stock, error)
	Upsert(ctx context.Context, s *Stock) error
	List(ctx context.Context) ([]*Stock, error)
}

type InMemoryReservationRepository struct {
	mu           sync.RWMutex
	reservations map[string]*InventoryReservation
}

func NewInMemoryReservationRepository() *InMemoryReservationRepository {
	return &InMemoryReservationRepository{reservations: make(map[string]*InventoryReservation)}
}

func (r *InMemoryReservationRepository) Create(_ context.Context, res *InventoryReservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reservations[res.ID] = res
	return nil
}

func (r *InMemoryReservationRepository) GetByID(_ context.Context, id string) (*InventoryReservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res, ok := r.reservations[id]
	if !ok {
		return nil, ErrReservationNotFound
	}
	return res, nil
}

func (r *InMemoryReservationRepository) Update(_ context.Context, res *InventoryReservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.reservations[res.ID]; !ok {
		return ErrReservationNotFound
	}
	r.reservations[res.ID] = res
	return nil
}

func (r *InMemoryReservationRepository) GetByOrderID(_ context.Context, orderID string) ([]*InventoryReservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*InventoryReservation
	for _, res := range r.reservations {
		if res.OrderID == orderID {
			result = append(result, res)
		}
	}
	return result, nil
}

var _ ReservationRepository = (*InMemoryReservationRepository)(nil)

type InMemoryStockRepository struct {
	mu     sync.RWMutex
	stocks map[string]*Stock // key: "productID:warehouseID"
}

func NewInMemoryStockRepository() *InMemoryStockRepository {
	return &InMemoryStockRepository{stocks: make(map[string]*Stock)}
}

func (r *InMemoryStockRepository) key(productID, warehouseID string) string {
	return productID + ":" + warehouseID
}

func (r *InMemoryStockRepository) Get(_ context.Context, productID, warehouseID string) (*Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.stocks[r.key(productID, warehouseID)]
	if !ok {
		return nil, ErrStockNotFound
	}
	return s, nil
}

func (r *InMemoryStockRepository) Upsert(_ context.Context, s *Stock) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stocks[r.key(s.ProductID, s.WarehouseID)] = s
	return nil
}

func (r *InMemoryStockRepository) List(_ context.Context) ([]*Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Stock
	for _, s := range r.stocks {
		result = append(result, s)
	}
	return result, nil
}

var _ StockRepository = (*InMemoryStockRepository)(nil)
