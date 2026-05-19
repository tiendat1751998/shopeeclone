package transcoding

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateJob(ctx context.Context, streamID string, inputURL string, profiles []VideoProfile) (*TranscodeJob, error) {
	if streamID == "" || inputURL == "" {
		return nil, ErrInvalidJobData
	}
	if len(profiles) == 0 {
		profiles = []VideoProfile{Profile480p, Profile720p}
	}
	for _, p := range profiles {
		if _, ok := SupportedProfiles[p]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedProfile, p)
		}
	}

	job := &TranscodeJob{
		ID:        uuid.New().String(),
		StreamID:  streamID,
		InputURL:  inputURL,
		Profiles:  profiles,
		Status:    JobPending,
		Progress:  0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.SaveJob(ctx, job); err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}
	return job, nil
}

func (s *Service) GetJobStatus(ctx context.Context, jobID string) (*TranscodeJob, error) {
	return s.repo.GetJob(ctx, jobID)
}

func (s *Service) ListJobs(ctx context.Context, streamID string) ([]*TranscodeJob, error) {
	if streamID != "" {
		return s.repo.ListJobsByStream(ctx, streamID)
	}
	return s.repo.ListJobs(ctx)
}

func (s *Service) CancelJob(ctx context.Context, jobID string) error {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return err
	}

	switch job.Status {
	case JobPending, JobProcessing:
		job.Status = JobFailed
		job.Error = "cancelled"
		job.UpdatedAt = time.Now().UTC()
		return s.repo.SaveJob(ctx, job)
	default:
		return ErrJobAlreadyStarted
	}
}

func (s *Service) StartJob(ctx context.Context, jobID string) error {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return err
	}
	if job.Status != JobPending {
		return ErrJobAlreadyStarted
	}
	now := time.Now().UTC()
	job.Status = JobProcessing
	job.StartedAt = &now
	job.UpdatedAt = now
	return s.repo.SaveJob(ctx, job)
}

func (s *Service) CompleteJob(ctx context.Context, jobID string, outputs []Output) error {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	job.Status = JobCompleted
	job.Outputs = outputs
	job.Progress = 100
	job.CompletedAt = &now
	job.UpdatedAt = now
	return s.repo.SaveJob(ctx, job)
}

func (s *Service) FailJob(ctx context.Context, jobID string, errMsg string) error {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	job.Status = JobFailed
	job.Error = errMsg
	job.CompletedAt = &now
	job.UpdatedAt = now
	return s.repo.SaveJob(ctx, job)
}
