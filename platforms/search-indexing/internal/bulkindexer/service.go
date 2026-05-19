package bulkindexer

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateJob(ctx context.Context, indexName string, totalDocuments int) (*BulkJob, error)
	SubmitBatch(ctx context.Context, jobID string, documents []map[string]interface{}, batchNumber int) (*DocumentBatch, error)
	MarkBatchProcessed(ctx context.Context, batchID string, failedCount int, errors []string) error
	GetJobProgress(ctx context.Context, id string) (*BulkJob, error)
	ListJobs(ctx context.Context) ([]*BulkJob, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateJob(ctx context.Context, indexName string, totalDocuments int) (*BulkJob, error) {
	job := &BulkJob{
		ID:             uuid.New().String(),
		IndexName:      indexName,
		TotalDocuments: totalDocuments,
		Status:         JobStatusPending,
		CreatedAt:      time.Now(),
	}
	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *service) SubmitBatch(ctx context.Context, jobID string, documents []map[string]interface{}, batchNumber int) (*DocumentBatch, error) {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job.Status == JobStatusCompleted {
		return nil, ErrJobCompleted
	}

	batch := &DocumentBatch{
		ID:          uuid.New().String(),
		JobID:       jobID,
		Documents:   documents,
		BatchNumber: batchNumber,
		Status:      BatchStatusPending,
	}
	if err := s.repo.CreateBatch(ctx, batch); err != nil {
		return nil, err
	}

	job.Status = JobStatusRunning
	s.repo.UpdateJob(ctx, job)

	return batch, nil
}

func (s *service) MarkBatchProcessed(ctx context.Context, batchID string, failedCount int, errors []string) error {
	batch, err := s.repo.GetBatch(ctx, batchID)
	if err != nil {
		return err
	}
	if failedCount > 0 {
		batch.Status = BatchStatusFailed
	} else {
		batch.Status = BatchStatusProcessed
	}
	s.repo.UpdateBatch(ctx, batch)

	job, err := s.repo.GetJob(ctx, batch.JobID)
	if err != nil {
		return err
	}

	job.ProcessedCount += len(batch.Documents) - failedCount
	job.FailedCount += failedCount
	if errors != nil {
		job.Errors = append(job.Errors, errors...)
	}

	if job.ProcessedCount+job.FailedCount >= job.TotalDocuments {
		job.Status = JobStatusCompleted
		now := time.Now()
		job.CompletedAt = &now
	} else {
		job.Status = JobStatusRunning
	}
	return s.repo.UpdateJob(ctx, job)
}

func (s *service) GetJobProgress(ctx context.Context, id string) (*BulkJob, error) {
	return s.repo.GetJob(ctx, id)
}

func (s *service) ListJobs(ctx context.Context) ([]*BulkJob, error) {
	return s.repo.ListJobs(ctx)
}
