package discovery

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("service instance not found")
)

type Watcher struct {
	ID      string
	Updates chan []*ServiceInstance
	Close   func()
}

type Service struct {
	repo              Repository
	mu                sync.RWMutex
	writers           map[string]*Watcher
	missedHeartbeats  map[string]int
	stopCh            chan struct{}
}

func NewService(repo Repository) *Service {
	s := &Service{
		repo:             repo,
		writers:          make(map[string]*Watcher),
		missedHeartbeats: make(map[string]int),
		stopCh:           make(chan struct{}),
	}
	go s.healthCheckLoop()
	return s
}

func (s *Service) Stop() {
	close(s.stopCh)
}

func (s *Service) Register(ctx context.Context, inst *ServiceInstance) error {
	return s.repo.Register(ctx, inst)
}

func (s *Service) Deregister(ctx context.Context, id string) error {
	return s.repo.Deregister(ctx, id)
}

func (s *Service) Heartbeat(ctx context.Context, id string) error {
	s.mu.Lock()
	s.missedHeartbeats[id] = 0
	s.mu.Unlock()
	return s.repo.Heartbeat(ctx, id)
}

func (s *Service) Discover(ctx context.Context, name string, region, zone string) ([]*ServiceInstance, error) {
	return s.repo.Discover(ctx, name, region, zone)
}

func (s *Service) ListServices(ctx context.Context) ([]*ServiceInstance, error) {
	return s.repo.ListServices(ctx)
}

func (s *Service) Watch(ctx context.Context, name string) *Watcher {
	w := &Watcher{
		ID:      name,
		Updates: make(chan []*ServiceInstance, 10),
	}
	s.mu.Lock()
	s.writers[name] = w
	s.mu.Unlock()

	instances, _ := s.repo.Discover(ctx, name, "", "")
	w.Updates <- instances

	w.Close = func() {
		s.mu.Lock()
		delete(s.writers, name)
		s.mu.Unlock()
	}
	return w
}

func (s *Service) healthCheckLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			instances, err := s.repo.ListServices(ctx)
			cancel()
			if err != nil {
				zap.L().Error("health check: failed to list services", zap.Error(err))
				continue
			}
			s.mu.Lock()
			for _, inst := range instances {
				if time.Since(inst.LastHeartbeat) > 15*time.Second {
					s.missedHeartbeats[inst.ID]++
					if s.missedHeartbeats[inst.ID] >= 3 {
						inst.Status = StatusDown
					}
				} else {
					s.missedHeartbeats[inst.ID] = 0
				}
			}
			s.mu.Unlock()
		case <-s.stopCh:
			return
		}
	}
}
