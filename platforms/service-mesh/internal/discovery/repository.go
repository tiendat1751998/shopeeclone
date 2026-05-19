package discovery

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Register(ctx context.Context, inst *ServiceInstance) error
	Deregister(ctx context.Context, id string) error
	Heartbeat(ctx context.Context, id string) error
	Discover(ctx context.Context, name string, region, zone string) ([]*ServiceInstance, error)
	ListServices(ctx context.Context) ([]*ServiceInstance, error)
	GetByID(ctx context.Context, id string) (*ServiceInstance, error)
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

func (r *InMemoryRepository) Register(ctx context.Context, inst *ServiceInstance) error {
	if inst.ID == "" {
		inst.ID = uuid.New().String()
	}
	inst.LastHeartbeat = time.Now()
	if inst.Status == "" {
		inst.Status = StatusUp
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[inst.ID] = inst
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
		return nil
	}
	inst.LastHeartbeat = time.Now()
	if inst.Status == StatusDown {
		inst.Status = StatusUp
	}
	return nil
}

func (r *InMemoryRepository) Discover(ctx context.Context, name string, region, zone string) ([]*ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*ServiceInstance
	for _, inst := range r.instances {
		if inst.Name != name {
			continue
		}
		if inst.Status != StatusUp {
			continue
		}
		if region != "" && inst.Region != region {
			continue
		}
		if zone != "" && inst.Zone != zone {
			continue
		}
		result = append(result, inst)
	}
	return result, nil
}

func (r *InMemoryRepository) ListServices(ctx context.Context) ([]*ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*ServiceInstance
	for _, inst := range r.instances {
		result = append(result, inst)
	}
	return result, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	inst, ok := r.instances[id]
	if !ok {
		return nil, nil
	}
	return inst, nil
}
