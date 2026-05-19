package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Process(ctx context.Context, orderID, userID string, amount, fee int64, currency string, method PaymentMethod) (*Payment, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if method != MethodCreditCard && method != MethodBankTransfer && method != MethodWallet && method != MethodCOD {
		return nil, ErrInvalidMethod
	}

	netAmount := amount - fee
	if netAmount < 0 {
		netAmount = 0
	}

	p := &Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		Method:        method,
		Status:        StatusPending,
		Fee:           fee,
		NetAmount:     netAmount,
		TransactionID: uuid.New().String(),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) Authorize(ctx context.Context, paymentID string) (*Payment, error) {
	p, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p.Status != StatusPending {
		return nil, ErrInvalidStatus
	}
	p.Status = StatusProcessing
	p.GatewayResponse = "authorized"
	p.AuthorizedAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Capture(ctx context.Context, paymentID string) (*Payment, error) {
	p, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p.Status != StatusProcessing {
		return nil, ErrInvalidStatus
	}
	p.Status = StatusCompleted
	p.GatewayResponse = "captured"
	p.CapturedAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Settle(ctx context.Context, paymentID string) (*Payment, error) {
	p, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p.Status != StatusCompleted {
		return nil, ErrInvalidStatus
	}
	p.GatewayResponse = "settled"
	p.SettledAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Refund(ctx context.Context, paymentID string) (*Payment, error) {
	p, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p.Status == StatusRefunded {
		return nil, ErrAlreadyRefunded
	}
	if p.Status != StatusCompleted && p.Status != StatusPartiallyRefunded {
		return nil, ErrInvalidStatus
	}
	p.Status = StatusRefunded
	p.GatewayResponse = "refunded"
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) PartialRefund(ctx context.Context, paymentID string, refundAmount int64) (*Payment, error) {
	p, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p.Status != StatusCompleted {
		return nil, ErrInvalidStatus
	}
	if refundAmount <= 0 || refundAmount > p.Amount {
		return nil, ErrRefundExceeds
	}
	p.Status = StatusPartiallyRefunded
	p.NetAmount = p.NetAmount - refundAmount
	if p.NetAmount < 0 {
		p.NetAmount = 0
	}
	p.GatewayResponse = "partially_refunded"
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Payment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]*Payment, int64, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *Service) GetByOrder(ctx context.Context, orderID string) ([]*Payment, error) {
	return s.repo.GetByOrder(ctx, orderID)
}
