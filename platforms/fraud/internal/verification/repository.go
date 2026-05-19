package verification

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, req *VerificationRequest) error
	Get(ctx context.Context, id string) (*VerificationRequest, error)
	GetByUserAndMethod(ctx context.Context, userID string, method VerificationMethod) (*VerificationRequest, error)
	Update(ctx context.Context, req *VerificationRequest) error
	GetKYCStatus(ctx context.Context, userID string) (*KYCStatus, error)
	SetKYCStatus(ctx context.Context, status *KYCStatus) error
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	requests    map[string]*VerificationRequest
	userReqs    map[string]string
	kycStatuses map[string]*KYCStatus
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		requests:    make(map[string]*VerificationRequest),
		userReqs:    make(map[string]string),
		kycStatuses: make(map[string]*KYCStatus),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, req *VerificationRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}
	req.Status = StatusPending
	req.MaxAttempts = 3
	r.requests[req.ID] = req
	r.userReqs[req.UserID+":"+string(req.Method)] = req.ID
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*VerificationRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	req, ok := r.requests[id]
	if !ok {
		return nil, ErrVerificationNotFound
	}
	return req, nil
}

func (r *InMemoryRepository) GetByUserAndMethod(ctx context.Context, userID string, method VerificationMethod) (*VerificationRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.userReqs[userID+":"+string(method)]
	if !ok {
		return nil, ErrVerificationNotFound
	}
	req, ok := r.requests[id]
	if !ok {
		return nil, ErrVerificationNotFound
	}
	return req, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, req *VerificationRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests[req.ID] = req
	return nil
}

func (r *InMemoryRepository) GetKYCStatus(ctx context.Context, userID string) (*KYCStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	status, ok := r.kycStatuses[userID]
	if !ok {
		return &KYCStatus{UserID: userID, IsVerified: false}, nil
	}
	return status, nil
}

func (r *InMemoryRepository) SetKYCStatus(ctx context.Context, status *KYCStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.kycStatuses[status.UserID] = status
	return nil
}
