package dispatcher

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, job *DispatchJob) error
	GetByID(ctx context.Context, id string) (*DispatchJob, error)
	GetPending(ctx context.Context) ([]*DispatchJob, error)
	GetFailed(ctx context.Context) ([]*DispatchJob, error)
	UpdateStatus(ctx context.Context, id, status, errMsg string) error
	IncrementRetry(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	jobs map[string]*DispatchJob
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		jobs: make(map[string]*DispatchJob),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, job *DispatchJob) error {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	if job.Status == "" {
		job.Status = "pending"
	}
	if job.MaxRetries == 0 {
		job.MaxRetries = 3
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*DispatchJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	j, ok := r.jobs[id]
	if !ok {
		return nil, nil
	}
	return j, nil
}

func (r *InMemoryRepository) GetPending(ctx context.Context) ([]*DispatchJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*DispatchJob
	for _, j := range r.jobs {
		if j.Status == "pending" {
			result = append(result, j)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) GetFailed(ctx context.Context) ([]*DispatchJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*DispatchJob
	for _, j := range r.jobs {
		if j.Status == "failed" {
			result = append(result, j)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) UpdateStatus(ctx context.Context, id, status, errMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	j, ok := r.jobs[id]
	if !ok {
		return nil
	}
	j.Status = status
	j.Error = errMsg
	j.UpdatedAt = time.Now()
	return nil
}

func (r *InMemoryRepository) IncrementRetry(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	j, ok := r.jobs[id]
	if !ok {
		return nil
	}
	j.RetryCount++
	j.UpdatedAt = time.Now()
	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.jobs, id)
	return nil
}
