package returns

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RequestReturn(ctx context.Context, r *Return) error {
	if r.ID == "" || r.OrderID == "" || r.UserID == "" {
		return ErrInvalidReturnData
	}
	if len(r.Items) == 0 {
		return fmt.Errorf("%w: at least one item required", ErrInvalidReturnData)
	}
	r.Status = ReturnStatusRequested
	r.RMNumber = fmt.Sprintf("RMA-%s-%d", r.ID, time.Now().Unix())
	r.CreatedAt = time.Now().UTC()
	return s.repo.Create(ctx, r)
}

func (s *Service) ApproveReturn(ctx context.Context, id string) error {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !IsValidReturnTransition(r.Status, ReturnStatusApproved) {
		return ErrInvalidReturnStatus
	}
	r.Status = ReturnStatusApproved
	now := time.Now().UTC()
	r.ApprovedAt = &now
	return s.repo.Update(ctx, r)
}

func (s *Service) RejectReturn(ctx context.Context, id string) error {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !IsValidReturnTransition(r.Status, ReturnStatusRejected) {
		return ErrInvalidReturnStatus
	}
	r.Status = ReturnStatusRejected
	return s.repo.Update(ctx, r)
}

func (s *Service) ReceiveReturn(ctx context.Context, id string) error {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !IsValidReturnTransition(r.Status, ReturnStatusReceived) {
		return ErrInvalidReturnStatus
	}
	r.Status = ReturnStatusReceived
	now := time.Now().UTC()
	r.ReceivedAt = &now
	return s.repo.Update(ctx, r)
}

func (s *Service) ProcessRefund(ctx context.Context, id string, refundAmount float64) error {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !IsValidReturnTransition(r.Status, ReturnStatusRefunded) {
		return ErrInvalidReturnStatus
	}
	r.Status = ReturnStatusRefunded
	r.RefundAmount = refundAmount
	now := time.Now().UTC()
	r.RefundedAt = &now
	return s.repo.Update(ctx, r)
}
