package payout

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	payoutRepo PayoutRepository
	batchRepo  BatchRepository
}

func NewService(payoutRepo PayoutRepository, batchRepo BatchRepository) *Service {
	return &Service{
		payoutRepo: payoutRepo,
		batchRepo:  batchRepo,
	}
}

func (s *Service) CreatePayout(ctx context.Context, sellerID string, amount, fee int64, paymentMethod PaymentMethod, periodStart, periodEnd string) (*Payout, error) {
	if amount <= 0 {
		return nil, ErrInvalidPayoutAmount
	}
	netAmount := amount - fee
	if netAmount < 0 {
		netAmount = 0
	}

	p := &Payout{
		ID:            uuid.New().String(),
		SellerID:      sellerID,
		Amount:        amount,
		Fee:           fee,
		NetAmount:     netAmount,
		PaymentMethod: paymentMethod,
		Status:        StatusPending,
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
	}

	if err := s.payoutRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ProcessPayout(ctx context.Context, payoutID string) (*Payout, error) {
	p, err := s.payoutRepo.GetByID(ctx, payoutID)
	if err != nil {
		return nil, err
	}
	if p.Status != StatusPending {
		return nil, ErrInvalidPayoutStatus
	}
	p.Status = StatusProcessing
	if err := s.payoutRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	p.Status = StatusCompleted
	now := time.Now().UTC().Format(time.RFC3339)
	p.CompletedAt = now
	if err := s.payoutRepo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) BatchPayout(ctx context.Context, payouts []*Payout) (*PayoutBatch, error) {
	totalAmount := int64(0)
	for _, p := range payouts {
		if p.Amount <= 0 {
			return nil, ErrInvalidPayoutAmount
		}
		totalAmount += p.NetAmount
	}

	batch := &PayoutBatch{
		ID:          uuid.New().String(),
		TotalAmount: totalAmount,
		Count:       len(payouts),
		Status:      StatusPending,
	}

	if err := s.batchRepo.Create(ctx, batch); err != nil {
		return nil, err
	}

	for _, p := range payouts {
		p.BatchID = batch.ID
		p.Status = StatusProcessing
		if err := s.payoutRepo.Update(ctx, p); err != nil {
			return nil, err
		}
	}

	batch.Status = StatusCompleted
	if err := s.batchRepo.Update(ctx, batch); err != nil {
		return nil, err
	}

	for _, p := range payouts {
		p.Status = StatusCompleted
		now := time.Now().UTC().Format(time.RFC3339)
		p.CompletedAt = now
		if err := s.payoutRepo.Update(ctx, p); err != nil {
			return nil, err
		}
	}

	return batch, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Payout, error) {
	return s.payoutRepo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]*Payout, int64, error) {
	return s.payoutRepo.List(ctx, offset, limit)
}
