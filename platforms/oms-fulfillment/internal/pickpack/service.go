package pickpack

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	pickListRepo  PickListRepository
	packingRepo   PackingRepository
	shipmentRepo  ShipmentRepository
}

func NewService(
	pickListRepo PickListRepository,
	packingRepo PackingRepository,
	shipmentRepo ShipmentRepository,
) *Service {
	return &Service{
		pickListRepo: pickListRepo,
		packingRepo:  packingRepo,
		shipmentRepo: shipmentRepo,
	}
}

func (s *Service) CreatePickList(ctx context.Context, pl *PickList) error {
	if pl.ID == "" || pl.OrderID == "" {
		return ErrInvalidPickData
	}
	pl.Status = PickListPending
	pl.CreatedAt = time.Now().UTC()
	return s.pickListRepo.Create(ctx, pl)
}

func (s *Service) CompletePick(ctx context.Context, id string) error {
	pl, err := s.pickListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if pl.Status != PickListInProgress && pl.Status != PickListPending {
		return fmt.Errorf("cannot complete pick list with status %s", pl.Status)
	}
	pl.Status = PickListCompleted
	now := time.Now().UTC()
	pl.CompletedAt = &now
	return s.pickListRepo.Update(ctx, pl)
}

func (s *Service) AssignPickList(ctx context.Context, id, assignee string) error {
	pl, err := s.pickListRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	pl.AssignedTo = assignee
	pl.Status = PickListInProgress
	return s.pickListRepo.Update(ctx, pl)
}

func (s *Service) CreatePacking(ctx context.Context, p *Packing) error {
	if p.ID == "" || p.PickListID == "" {
		return ErrInvalidPackData
	}
	p.Status = PackingPending
	p.CreatedAt = time.Now().UTC()
	return s.packingRepo.Create(ctx, p)
}

func (s *Service) CreateShipment(ctx context.Context, sh *Shipment) error {
	if sh.ID == "" || sh.PackingID == "" {
		return ErrInvalidShipData
	}
	sh.Status = ShipmentPending
	sh.CreatedAt = time.Now().UTC()
	return s.shipmentRepo.Create(ctx, sh)
}
