package pickpack

import (
	"context"
	"sync"
)

type PickListRepository interface {
	Create(ctx context.Context, pl *PickList) error
	GetByID(ctx context.Context, id string) (*PickList, error)
	Update(ctx context.Context, pl *PickList) error
}

type PackingRepository interface {
	Create(ctx context.Context, p *Packing) error
	GetByID(ctx context.Context, id string) (*Packing, error)
}

type ShipmentRepository interface {
	Create(ctx context.Context, s *Shipment) error
	GetByID(ctx context.Context, id string) (*Shipment, error)
}

type InMemoryPickListRepository struct {
	mu       sync.RWMutex
	lists    map[string]*PickList
}

func NewInMemoryPickListRepository() *InMemoryPickListRepository {
	return &InMemoryPickListRepository{lists: make(map[string]*PickList)}
}

func (r *InMemoryPickListRepository) Create(_ context.Context, pl *PickList) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lists[pl.ID] = pl
	return nil
}

func (r *InMemoryPickListRepository) GetByID(_ context.Context, id string) (*PickList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pl, ok := r.lists[id]
	if !ok {
		return nil, ErrPickListNotFound
	}
	return pl, nil
}

func (r *InMemoryPickListRepository) Update(_ context.Context, pl *PickList) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.lists[pl.ID]; !ok {
		return ErrPickListNotFound
	}
	r.lists[pl.ID] = pl
	return nil
}

var _ PickListRepository = (*InMemoryPickListRepository)(nil)

type InMemoryPackingRepository struct {
	mu      sync.RWMutex
	packings map[string]*Packing
}

func NewInMemoryPackingRepository() *InMemoryPackingRepository {
	return &InMemoryPackingRepository{packings: make(map[string]*Packing)}
}

func (r *InMemoryPackingRepository) Create(_ context.Context, p *Packing) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.packings[p.ID] = p
	return nil
}

func (r *InMemoryPackingRepository) GetByID(_ context.Context, id string) (*Packing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.packings[id]
	if !ok {
		return nil, ErrPackingNotFound
	}
	return p, nil
}

var _ PackingRepository = (*InMemoryPackingRepository)(nil)

type InMemoryShipmentRepository struct {
	mu        sync.RWMutex
	shipments map[string]*Shipment
}

func NewInMemoryShipmentRepository() *InMemoryShipmentRepository {
	return &InMemoryShipmentRepository{shipments: make(map[string]*Shipment)}
}

func (r *InMemoryShipmentRepository) Create(_ context.Context, s *Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shipments[s.ID] = s
	return nil
}

func (r *InMemoryShipmentRepository) GetByID(_ context.Context, id string) (*Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.shipments[id]
	if !ok {
		return nil, ErrShipmentNotFound
	}
	return s, nil
}

var _ ShipmentRepository = (*InMemoryShipmentRepository)(nil)
