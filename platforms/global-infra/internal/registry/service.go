package registry

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Repository interface {
	Register(ctx context.Context, instance *ServiceInstance) error
	Deregister(ctx context.Context, id string) error
	Heartbeat(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*ServiceInstance, error)
	Discover(ctx context.Context, name string, region string) ([]*ServiceInstance, error)
	List(ctx context.Context) ([]*ServiceInstance, error)
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	instances map[string]*ServiceInstance
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		instances: make(map[string]*ServiceInstance),
	}
}

func (r *InMemoryRepository) Register(ctx context.Context, instance *ServiceInstance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	instance.LastHeartbeat = now
	instance.RegisteredAt = now
	if instance.Status == "" {
		instance.Status = StatusUp
	}
	r.instances[instance.ID] = instance
	return nil
}

func (r *InMemoryRepository) Deregister(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.instances, id)
	return nil
}

func (r *InMemoryRepository) Heartbeat(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inst, ok := r.instances[id]
	if !ok {
		return fmt.Errorf("instance not found: %s", id)
	}
	inst.LastHeartbeat = time.Now()
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inst, ok := r.instances[id]
	if !ok {
		return nil, nil
	}
	return inst, nil
}

func (r *InMemoryRepository) Discover(ctx context.Context, name string, region string) ([]*ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*ServiceInstance
	for _, inst := range r.instances {
		if inst.Name != name {
			continue
		}
		if region != "" && inst.Region != region {
			continue
		}
		if inst.Status == StatusUp {
			result = append(result, inst)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*ServiceInstance, 0, len(r.instances))
	for _, inst := range r.instances {
		result = append(result, inst)
	}
	return result, nil
}

type HealthChecker struct {
	client *http.Client
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (hc *HealthChecker) Check(endpoint string) bool {
	if endpoint == "" {
		return false
	}
	resp, err := hc.client.Get(endpoint)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

type Service struct {
	repo         Repository
	healthChecker *HealthChecker
}

func NewService(repo Repository, hc *HealthChecker) *Service {
	return &Service{repo: repo, healthChecker: hc}
}

func (s *Service) Register(ctx context.Context, instance *ServiceInstance) (*ServiceInstance, error) {
	if instance.ID == "" {
		return nil, fmt.Errorf("id is required")
	}
	if instance.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if instance.Address == "" {
		return nil, fmt.Errorf("address is required")
	}
	if instance.Port <= 0 {
		return nil, fmt.Errorf("port is required")
	}
	if instance.HealthEndpoint == "" {
		instance.HealthEndpoint = fmt.Sprintf("http://%s:%d/health/live", instance.Address, instance.Port)
	}
	if err := s.repo.Register(ctx, instance); err != nil {
		return nil, err
	}
	return instance, nil
}

func (s *Service) Deregister(ctx context.Context, id string) error {
	return s.repo.Deregister(ctx, id)
}

func (s *Service) Heartbeat(ctx context.Context, id string) error {
	inst, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if inst == nil {
		return fmt.Errorf("instance not found: %s", id)
	}
	if s.healthChecker != nil && !s.healthChecker.Check(inst.HealthEndpoint) {
		return fmt.Errorf("health check failed for instance: %s", id)
	}
	return s.repo.Heartbeat(ctx, id)
}

func (s *Service) Discover(ctx context.Context, name string, region string) ([]*ServiceInstance, error) {
	return s.repo.Discover(ctx, name, region)
}

func (s *Service) List(ctx context.Context) ([]*ServiceInstance, error) {
	return s.repo.List(ctx)
}
