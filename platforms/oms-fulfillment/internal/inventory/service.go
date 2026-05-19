package inventory

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	reservationRepo ReservationRepository
	stockRepo       StockRepository
}

func NewService(reservationRepo ReservationRepository, stockRepo StockRepository) *Service {
	return &Service{
		reservationRepo: reservationRepo,
		stockRepo:       stockRepo,
	}
}

func (s *Service) Reserve(ctx context.Context, req ReserveRequest) (*InventoryReservation, error) {
	stock, err := s.stockRepo.Get(ctx, req.ProductID, "default")
	if err != nil {
		return nil, fmt.Errorf("check stock: %w", err)
	}
	if stock.Available < req.Quantity {
		return nil, ErrInsufficientStock
	}
	res := &InventoryReservation{
		ID:        fmt.Sprintf("res-%d", time.Now().UnixNano()),
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		SKU:       req.SKU,
		Quantity:  req.Quantity,
		Status:    ReservationReserved,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	stock.Available -= req.Quantity
	stock.Reserved += req.Quantity
	if err := s.stockRepo.Upsert(ctx, stock); err != nil {
		return nil, err
	}
	if err := s.reservationRepo.Create(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Release(ctx context.Context, reservationID string) error {
	res, err := s.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		return err
	}
	if res.Status != ReservationReserved {
		return nil
	}
	stock, err := s.stockRepo.Get(ctx, res.ProductID, "default")
	if err != nil {
		return err
	}
	stock.Available += res.Quantity
	stock.Reserved -= res.Quantity
	if err := s.stockRepo.Upsert(ctx, stock); err != nil {
		return err
	}
	res.Status = ReservationReleased
	res.UpdatedAt = time.Now().UTC()
	return s.reservationRepo.Update(ctx, res)
}

func (s *Service) Consume(ctx context.Context, reservationID string) error {
	res, err := s.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		return err
	}
	if res.Status != ReservationReserved {
		return nil
	}
	stock, err := s.stockRepo.Get(ctx, res.ProductID, "default")
	if err != nil {
		return err
	}
	stock.Reserved -= res.Quantity
	stock.Total -= res.Quantity
	if err := s.stockRepo.Upsert(ctx, stock); err != nil {
		return err
	}
	res.Status = ReservationConsumed
	res.UpdatedAt = time.Now().UTC()
	return s.reservationRepo.Update(ctx, res)
}

func (s *Service) CheckAvailability(ctx context.Context, productID string, quantity int) (bool, error) {
	stock, err := s.stockRepo.Get(ctx, productID, "default")
	if err != nil {
		return false, err
	}
	return stock.Available >= quantity, nil
}

func (s *Service) GetReservationsByOrder(ctx context.Context, orderID string) ([]*InventoryReservation, error) {
	return s.reservationRepo.GetByOrderID(ctx, orderID)
}

func (s *Service) GetStock(ctx context.Context, productID string) (*Stock, error) {
	return s.stockRepo.Get(ctx, productID, "default")
}

func (s *Service) ListStock(ctx context.Context) ([]*Stock, error) {
	return s.stockRepo.List(ctx)
}
