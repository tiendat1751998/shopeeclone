package warehouse

import (
	"context"
	"sync"
)

type WarehouseRepository interface {
	Create(ctx context.Context, w *Warehouse) error
	GetByID(ctx context.Context, id string) (*Warehouse, error)
	Update(ctx context.Context, w *Warehouse) error
	List(ctx context.Context) ([]*Warehouse, error)
}

type ZoneRepository interface {
	Create(ctx context.Context, z *Zone) error
	GetByWarehouse(ctx context.Context, warehouseID string) ([]*Zone, error)
}

type MovementRepository interface {
	Create(ctx context.Context, m *InventoryMovement) error
	List(ctx context.Context) ([]*InventoryMovement, error)
}

type InMemoryWarehouseRepository struct {
	mu         sync.RWMutex
	warehouses map[string]*Warehouse
}

func NewInMemoryWarehouseRepository() *InMemoryWarehouseRepository {
	return &InMemoryWarehouseRepository{warehouses: make(map[string]*Warehouse)}
}

func (r *InMemoryWarehouseRepository) Create(_ context.Context, w *Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.warehouses[w.ID] = w
	return nil
}

func (r *InMemoryWarehouseRepository) GetByID(_ context.Context, id string) (*Warehouse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.warehouses[id]
	if !ok {
		return nil, ErrWarehouseNotFound
	}
	return w, nil
}

func (r *InMemoryWarehouseRepository) Update(_ context.Context, w *Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.warehouses[w.ID]; !ok {
		return ErrWarehouseNotFound
	}
	r.warehouses[w.ID] = w
	return nil
}

func (r *InMemoryWarehouseRepository) List(_ context.Context) ([]*Warehouse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Warehouse
	for _, w := range r.warehouses {
		result = append(result, w)
	}
	return result, nil
}

var _ WarehouseRepository = (*InMemoryWarehouseRepository)(nil)

type InMemoryZoneRepository struct {
	mu    sync.RWMutex
	zones map[string]*Zone
}

func NewInMemoryZoneRepository() *InMemoryZoneRepository {
	return &InMemoryZoneRepository{zones: make(map[string]*Zone)}
}

func (r *InMemoryZoneRepository) Create(_ context.Context, z *Zone) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.zones[z.ID] = z
	return nil
}

func (r *InMemoryZoneRepository) GetByWarehouse(_ context.Context, warehouseID string) ([]*Zone, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Zone
	for _, z := range r.zones {
		if z.WarehouseID == warehouseID {
			result = append(result, z)
		}
	}
	return result, nil
}

var _ ZoneRepository = (*InMemoryZoneRepository)(nil)

type InMemoryMovementRepository struct {
	mu        sync.RWMutex
	movements map[string]*InventoryMovement
}

func NewInMemoryMovementRepository() *InMemoryMovementRepository {
	return &InMemoryMovementRepository{movements: make(map[string]*InventoryMovement)}
}

func (r *InMemoryMovementRepository) Create(_ context.Context, m *InventoryMovement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.movements[m.ID] = m
	return nil
}

func (r *InMemoryMovementRepository) List(_ context.Context) ([]*InventoryMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*InventoryMovement
	for _, m := range r.movements {
		result = append(result, m)
	}
	return result, nil
}

var _ MovementRepository = (*InMemoryMovementRepository)(nil)
