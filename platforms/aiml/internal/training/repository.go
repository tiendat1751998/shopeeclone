package training

import (
	"context"
	"sync"
)

type Repository interface {
	Store(ctx context.Context, job *TrainingJob) error
	Get(ctx context.Context, id string) (*TrainingJob, error)
	List(ctx context.Context) ([]*TrainingJob, error)
	Update(ctx context.Context, job *TrainingJob) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	jobs map[string]*TrainingJob
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		jobs: make(map[string]*TrainingJob),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, job *TrainingJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.jobs[job.ID]; ok {
		return ErrJobExists
	}
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*TrainingJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[id]
	if !ok {
		return nil, ErrJobNotFound
	}
	return j, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*TrainingJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*TrainingJob, 0, len(r.jobs))
	for _, j := range r.jobs {
		result = append(result, j)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, job *TrainingJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.jobs[job.ID]; !ok {
		return ErrJobNotFound
	}
	r.jobs[job.ID] = job
	return nil
}
