package warehouse

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	warehouseRepo WarehouseRepository
	zoneRepo      ZoneRepository
	movementRepo  MovementRepository
}

func NewService(
	warehouseRepo WarehouseRepository,
	zoneRepo ZoneRepository,
	movementRepo MovementRepository,
) *Service {
	return &Service{
		warehouseRepo: warehouseRepo,
		zoneRepo:      zoneRepo,
		movementRepo:  movementRepo,
	}
}

func (s *Service) CreateWarehouse(ctx context.Context, w *Warehouse) error {
	if w.ID == "" || w.Name == "" {
		return fmt.Errorf("%w: id and name required", ErrInvalidMovement)
	}
	return s.warehouseRepo.Create(ctx, w)
}

func (s *Service) ListWarehouses(ctx context.Context) ([]*Warehouse, error) {
	return s.warehouseRepo.List(ctx)
}

func (s *Service) GetWarehouse(ctx context.Context, id string) (*Warehouse, error) {
	return s.warehouseRepo.GetByID(ctx, id)
}

func (s *Service) CreateZone(ctx context.Context, z *Zone) error {
	if z.ID == "" || z.WarehouseID == "" {
		return fmt.Errorf("%w: id and warehouse_id required", ErrInvalidMovement)
	}
	return s.zoneRepo.Create(ctx, z)
}

func (s *Service) GetZones(ctx context.Context, warehouseID string) ([]*Zone, error) {
	return s.zoneRepo.GetByWarehouse(ctx, warehouseID)
}

func (s *Service) RecordMovement(ctx context.Context, m *InventoryMovement) error {
	if m.ID == "" || m.ProductID == "" || m.WarehouseID == "" {
		return ErrInvalidMovement
	}
	m.CreatedAt = time.Now().UTC()
	return s.movementRepo.Create(ctx, m)
}

func (s *Service) ListMovements(ctx context.Context) ([]*InventoryMovement, error) {
	return s.movementRepo.List(ctx)
}
