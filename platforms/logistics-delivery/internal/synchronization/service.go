package synchronization

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type SyncEvent struct {
	ID        string    `json:"id"`
	Entity    string    `json:"entity"`
	Action    string    `json:"action"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

type Handler func(ctx context.Context, event *SyncEvent) error

type Service struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

func NewService() *Service {
	return &Service{
		handlers: make(map[string]Handler),
	}
}

func (s *Service) RegisterHandler(entity, action string, handler Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := entity + ":" + action
	s.handlers[key] = handler
}

func (s *Service) ProcessEvent(ctx context.Context, event *SyncEvent) error {
	raw, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	event.Payload = raw
	s.mu.RLock()
	handler, ok := s.handlers[event.Entity+":"+event.Action]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no handler for %s:%s", event.Entity, event.Action)
	}
	return handler(ctx, event)
}

func (s *Service) Broadcast(ctx context.Context, events []*SyncEvent) error {
	for _, e := range events {
		if err := s.ProcessEvent(ctx, e); err != nil {
			return err
		}
	}
	return nil
}
