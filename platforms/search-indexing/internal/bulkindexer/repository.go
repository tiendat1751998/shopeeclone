package bulkindexer

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	CreateJob(ctx context.Context, job *BulkJob) error
	GetJob(ctx context.Context, id string) (*BulkJob, error)
	UpdateJob(ctx context.Context, job *BulkJob) error
	ListJobs(ctx context.Context) ([]*BulkJob, error)
	DeleteJob(ctx context.Context, id string) error
	CreateBatch(ctx context.Context, batch *DocumentBatch) error
	GetBatch(ctx context.Context, id string) (*DocumentBatch, error)
	UpdateBatch(ctx context.Context, batch *DocumentBatch) error
	ListBatchesByJob(ctx context.Context, jobID string) ([]*DocumentBatch, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	jobs    map[string]*BulkJob
	batches map[string]*DocumentBatch
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		jobs:    make(map[string]*BulkJob),
		batches: make(map[string]*DocumentBatch),
	}
}

func (r *InMemoryRepository) CreateJob(_ context.Context, job *BulkJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.jobs[job.ID]; ok {
		return ErrJobAlreadyExists
	}
	job.CreatedAt = time.Now()
	job.Status = JobStatusPending
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryRepository) GetJob(_ context.Context, id string) (*BulkJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.jobs[id]
	if !ok {
		return nil, ErrJobNotFound
	}
	return job, nil
}

func (r *InMemoryRepository) UpdateJob(_ context.Context, job *BulkJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.jobs[job.ID]; !ok {
		return ErrJobNotFound
	}
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryRepository) ListJobs(_ context.Context) ([]*BulkJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*BulkJob, 0, len(r.jobs))
	for _, j := range r.jobs {
		list = append(list, j)
	}
	return list, nil
}

func (r *InMemoryRepository) DeleteJob(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.jobs[id]; !ok {
		return ErrJobNotFound
	}
	delete(r.jobs, id)
	return nil
}

func (r *InMemoryRepository) CreateBatch(_ context.Context, batch *DocumentBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	batch.Status = BatchStatusPending
	r.batches[batch.ID] = batch
	return nil
}

func (r *InMemoryRepository) GetBatch(_ context.Context, id string) (*DocumentBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	batch, ok := r.batches[id]
	if !ok {
		return nil, ErrBatchNotFound
	}
	return batch, nil
}

func (r *InMemoryRepository) UpdateBatch(_ context.Context, batch *DocumentBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.batches[batch.ID]; !ok {
		return ErrBatchNotFound
	}
	r.batches[batch.ID] = batch
	return nil
}

func (r *InMemoryRepository) ListBatchesByJob(_ context.Context, jobID string) ([]*DocumentBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*DocumentBatch, 0)
	for _, b := range r.batches {
		if b.JobID == jobID {
			list = append(list, b)
		}
	}
	return list, nil
}
