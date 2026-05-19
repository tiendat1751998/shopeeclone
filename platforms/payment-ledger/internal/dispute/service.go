package dispute

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

func (s *Service) OpenDispute(ctx context.Context, transactionID, paymentID, userID, reason string, amount int64) (*Dispute, error) {
	d := &Dispute{
		ID:            uuid.New().String(),
		TransactionID: transactionID,
		PaymentID:     paymentID,
		UserID:        userID,
		Reason:        reason,
		Amount:        amount,
		Status:        StatusOpened,
		Evidence:      []string{},
	}

	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) SubmitEvidence(ctx context.Context, disputeID, evidenceItem string) (*Dispute, error) {
	d, err := s.repo.GetByID(ctx, disputeID)
	if err != nil {
		return nil, err
	}
	if d.Status == StatusResolved || d.Status == StatusClosed {
		return nil, ErrAlreadyResolved
	}
	d.Status = StatusUnderReview
	d.Evidence = append(d.Evidence, evidenceItem)
	if err := s.repo.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Resolve(ctx context.Context, disputeID string, resolution Resolution, notes string) (*Dispute, error) {
	d, err := s.repo.GetByID(ctx, disputeID)
	if err != nil {
		return nil, err
	}
	if d.Status == StatusClosed {
		return nil, ErrAlreadyResolved
	}
	d.Status = StatusResolved
	d.Resolution = resolution
	d.ResolvedAt = time.Now().UTC().Format(time.RFC3339)
	if notes != "" {
		d.Evidence = append(d.Evidence, "resolution: "+notes)
	}
	if err := s.repo.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Appeal(ctx context.Context, disputeID, reason string) (*Dispute, error) {
	d, err := s.repo.GetByID(ctx, disputeID)
	if err != nil {
		return nil, err
	}
	if d.Status != StatusResolved {
		return nil, ErrInvalidDisputeStatus
	}
	d.Status = StatusUnderReview
	d.Evidence = append(d.Evidence, "appeal: "+reason)
	if err := s.repo.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]*Dispute, int64, error) {
	return s.repo.List(ctx, offset, limit)
}
