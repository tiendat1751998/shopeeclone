package replay

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ReplayIDGenerator interface {
	Generate() string
}

type Service struct {
	mu           sync.RWMutex
	processedIDs map[string]time.Time
	ttl          time.Duration
}

func NewService(ttl time.Duration) *Service {
	return &Service{
		processedIDs: make(map[string]time.Time),
		ttl:          ttl,
	}
}

func (s *Service) IsProcessed(ctx context.Context, replayID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.processedIDs[replayID]
	if !ok {
		return false, nil
	}
	if time.Since(t) > s.ttl {
		return false, nil
	}
	return true, nil
}

func (s *Service) MarkProcessed(ctx context.Context, replayID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processedIDs[replayID] = time.Now().UTC()
	if len(s.processedIDs) > 10000 {
		s.cleanup()
	}
	return nil
}

func (s *Service) cleanup() {
	cutoff := time.Now().UTC().Add(-s.ttl)
	for id, t := range s.processedIDs {
		if t.Before(cutoff) {
			delete(s.processedIDs, id)
		}
	}
}

func (s *Service) ProcessWithReplayGuard(ctx context.Context, replayID string, fn func(context.Context) error) error {
	dup, err := s.IsProcessed(ctx, replayID)
	if err != nil {
		return fmt.Errorf("replay check: %w", err)
	}
	if dup {
		return nil
	}
	if err := fn(ctx); err != nil {
		return err
	}
	return s.MarkProcessed(ctx, replayID)
}
