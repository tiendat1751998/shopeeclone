package verification

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo            Repository
	codeExpiry      time.Duration
	maxAttempts     int
}

func NewService(repo Repository, codeExpiryMinutes int) *Service {
	return &Service{
		repo:        repo,
		codeExpiry:  time.Duration(codeExpiryMinutes) * time.Minute,
		maxAttempts: 3,
	}
}

func (s *Service) InitiateVerification(ctx context.Context, userID string, method VerificationMethod, target string) (*VerificationRequest, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	now := time.Now()

	req := &VerificationRequest{
		ID:          uuid.New().String(),
		UserID:      userID,
		Method:      method,
		Target:      target,
		Code:        code,
		Status:      StatusPending,
		MaxAttempts: s.maxAttempts,
		CreatedAt:   now,
		ExpiresAt:   now.Add(s.codeExpiry),
	}

	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *Service) VerifyCode(ctx context.Context, requestID, code string) (*VerificationRequest, error) {
	req, err := s.repo.Get(ctx, requestID)
	if err != nil {
		return nil, err
	}

	if req.Status == StatusVerified {
		return req, nil
	}

	if time.Now().After(req.ExpiresAt) {
		req.Status = StatusExpired
		s.repo.Update(ctx, req)
		return nil, ErrCodeExpired
	}

	req.Attempts++
	if req.Attempts > req.MaxAttempts {
		req.Status = StatusFailed
		s.repo.Update(ctx, req)
		return nil, ErrMaxAttempts
	}

	if req.Code != code {
		s.repo.Update(ctx, req)
		return nil, ErrCodeMismatch
	}

	now := time.Now()
	req.Status = StatusVerified
	req.VerifiedAt = &now
	s.repo.Update(ctx, req)

	return req, nil
}

func (s *Service) CheckKYCStatus(ctx context.Context, userID string) (*KYCStatus, error) {
	return s.repo.GetKYCStatus(ctx, userID)
}

func (s *Service) SetKYCStatus(ctx context.Context, status *KYCStatus) error {
	return s.repo.SetKYCStatus(ctx, status)
}
