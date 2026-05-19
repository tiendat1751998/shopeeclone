package blacklist

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

func (s *Service) Check(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {
	resp := &CheckResponse{Blocked: false}

	checks := map[BlacklistType]string{
		BlacklistUser:   req.UserID,
		BlacklistIP:     req.IP,
		BlacklistDevice: req.DeviceID,
	}

	if req.CardNumber != "" {
		checks[BlacklistCard] = req.CardNumber
	}

	for bt, val := range checks {
		if val == "" {
			continue
		}
		entry, err := s.repo.GetByTypeAndValue(ctx, bt, val)
		if err == nil && entry != nil && entry.IsActive {
			resp.Blocked = true
			resp.Reasons = append(resp.Reasons, entry.Reason)
			resp.Entries = append(resp.Entries, *entry)
		}
	}

	return resp, nil
}

func (s *Service) Add(ctx context.Context, entry *BlacklistEntry) error {
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	entry.IsActive = true
	return s.repo.Add(ctx, entry)
}

func (s *Service) Remove(ctx context.Context, id string) error {
	return s.repo.Remove(ctx, id)
}

func (s *Service) RemoveByValue(ctx context.Context, bt BlacklistType, value string) error {
	entry, err := s.repo.GetByTypeAndValue(ctx, bt, value)
	if err != nil {
		return err
	}
	return s.repo.Remove(ctx, entry.ID)
}

func (s *Service) BulkImport(ctx context.Context, entries []BlacklistEntry) (int, error) {
	added := 0
	for i := range entries {
		if entries[i].ID == "" {
			entries[i].ID = uuid.New().String()
		}
		if entries[i].CreatedAt.IsZero() {
			entries[i].CreatedAt = time.Now()
		}
		entries[i].IsActive = true
		if err := s.repo.Add(ctx, &entries[i]); err != nil {
			return added, fmt.Errorf("failed to add entry at index %d: %w", i, err)
		}
		added++
	}
	return added, nil
}

func (s *Service) ExpireEntries(ctx context.Context) (int, error) {
	expired := 0
	now := time.Now()

	list, err := s.repo.ListAll(ctx)
	if err != nil {
		return 0, err
	}

	for _, entry := range list {
		if entry.ExpiresAt != nil && now.After(*entry.ExpiresAt) && entry.IsActive {
			entry.IsActive = false
			if err := s.repo.Update(ctx, entry); err != nil {
				return expired, err
			}
			expired++
		}
	}
	return expired, nil
}
