package push

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	RegisterDevice(ctx context.Context, d *PushDevice) error
	GetDevice(ctx context.Context, id string) (*PushDevice, error)
	ListDevicesByUser(ctx context.Context, userID string) ([]*PushDevice, error)
	ListActiveDevicesByUser(ctx context.Context, userID string) ([]*PushDevice, error)
	MarkDeviceInactive(ctx context.Context, id string) error
	DeleteDevice(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	devices map[string]*PushDevice
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		devices: make(map[string]*PushDevice),
	}
}

func (r *InMemoryRepository) RegisterDevice(ctx context.Context, d *PushDevice) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	d.Active = true
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.devices {
		if existing.UserID == d.UserID && existing.Token == d.Token {
			existing.UpdatedAt = now
			existing.Active = true
			*d = *existing
			return nil
		}
	}

	r.devices[d.ID] = d
	return nil
}

func (r *InMemoryRepository) GetDevice(ctx context.Context, id string) (*PushDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.devices[id]
	if !ok {
		return nil, nil
	}
	return d, nil
}

func (r *InMemoryRepository) ListDevicesByUser(ctx context.Context, userID string) ([]*PushDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*PushDevice
	for _, d := range r.devices {
		if d.UserID == userID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) ListActiveDevicesByUser(ctx context.Context, userID string) ([]*PushDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*PushDevice
	for _, d := range r.devices {
		if d.UserID == userID && d.Active {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) MarkDeviceInactive(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	d, ok := r.devices[id]
	if !ok {
		return nil
	}
	d.Active = false
	d.UpdatedAt = time.Now()
	return nil
}

func (r *InMemoryRepository) DeleteDevice(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.devices, id)
	return nil
}
