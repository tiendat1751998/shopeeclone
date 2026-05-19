package ordermanagement

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, o *Order) error {
	if o.ID == "" || o.UserID == "" {
		return ErrInvalidOrderData
	}
	if len(o.Items) == 0 {
		return ErrEmptyItems
	}
	if o.Status == "" {
		o.Status = OrderStatusPending
	}
	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	return s.repo.Create(ctx, o)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter OrderFilter) ([]*Order, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return s.repo.List(ctx, filter)
}

func (s *Service) UpdateStatus(ctx context.Context, id string, newStatus OrderStatus) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !IsValidOrderTransition(order.Status, newStatus) {
		return ErrInvalidStatusTransition
	}
	order.Status = newStatus
	order.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, order)
}

func (s *Service) Cancel(ctx context.Context, id string) error {
	return s.UpdateStatus(ctx, id, OrderStatusCancelled)
}
