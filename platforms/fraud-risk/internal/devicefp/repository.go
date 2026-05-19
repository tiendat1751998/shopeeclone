package devicefp

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	Save(ctx context.Context, profile *DeviceProfile, hash string) error
	Get(ctx context.Context, deviceID string) (*DeviceProfile, error)
	List(ctx context.Context) ([]*DeviceProfile, error)
	GetByHash(ctx context.Context, hash string) (*DeviceProfile, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	profiles map[string]*DeviceProfile
	hashIdx  map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		profiles: make(map[string]*DeviceProfile),
		hashIdx:  make(map[string]string),
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, profile *DeviceProfile, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profiles[profile.DeviceID] = profile
	if hash != "" {
		r.hashIdx[hash] = profile.DeviceID
	}
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, deviceID string) (*DeviceProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	profile, ok := r.profiles[deviceID]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	return profile, nil
}

func (r *InMemoryRepository) GetByHash(ctx context.Context, hash string) (*DeviceProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	deviceID, ok := r.hashIdx[hash]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	profile, ok := r.profiles[deviceID]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	return profile, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*DeviceProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*DeviceProfile, 0, len(r.profiles))
	for _, p := range r.profiles {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryRepository) UpdateLastSeen(ctx context.Context, deviceID string, seen time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	profile, ok := r.profiles[deviceID]
	if !ok {
		return ErrDeviceNotFound
	}
	profile.LastSeen = seen
	return nil
}
