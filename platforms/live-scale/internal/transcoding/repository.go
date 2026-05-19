package transcoding

import (
	"context"
	"sync"
)

type Repository interface {
	SaveJob(ctx context.Context, job *TranscodeJob) error
	GetJob(ctx context.Context, id string) (*TranscodeJob, error)
	ListJobs(ctx context.Context) ([]*TranscodeJob, error)
	ListJobsByStream(ctx context.Context, streamID string) ([]*TranscodeJob, error)
	DeleteJob(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	jobs map[string]*TranscodeJob
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		jobs: make(map[string]*TranscodeJob),
	}
}

func (r *InMemoryRepository) SaveJob(_ context.Context, job *TranscodeJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryRepository) GetJob(_ context.Context, id string) (*TranscodeJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.jobs[id]
	if !ok {
		return nil, ErrJobNotFound
	}
	return job, nil
}

func (r *InMemoryRepository) ListJobs(_ context.Context) ([]*TranscodeJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*TranscodeJob, 0, len(r.jobs))
	for _, job := range r.jobs {
		result = append(result, job)
	}
	return result, nil
}

func (r *InMemoryRepository) ListJobsByStream(_ context.Context, streamID string) ([]*TranscodeJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*TranscodeJob
	for _, job := range r.jobs {
		if job.StreamID == streamID {
			result = append(result, job)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) DeleteJob(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.jobs, id)
	return nil
}
