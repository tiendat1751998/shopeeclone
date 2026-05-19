package audience

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateSegment(ctx context.Context, s *Segment) error
	ListSegments(ctx context.Context) ([]*Segment, error)
	GetSegmentByID(ctx context.Context, id string) (*Segment, error)
	CreateUser(ctx context.Context, u *UserProfile) error
	GetUserByID(ctx context.Context, id string) (*UserProfile, error)
	ListUsers(ctx context.Context) ([]*UserProfile, error)
	AddToSegment(ctx context.Context, userID string, segmentID string) error
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	segments map[string]*Segment
	users    map[string]*UserProfile
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		segments: make(map[string]*Segment),
		users:    make(map[string]*UserProfile),
	}
}

func (r *InMemoryRepository) CreateSegment(ctx context.Context, s *Segment) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.CreatedAt = time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.segments[s.ID] = s
	return nil
}

func (r *InMemoryRepository) ListSegments(ctx context.Context) ([]*Segment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Segment
	for _, s := range r.segments {
		result = append(result, s)
	}
	return result, nil
}

func (r *InMemoryRepository) GetSegmentByID(ctx context.Context, id string) (*Segment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.segments[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *InMemoryRepository) CreateUser(ctx context.Context, u *UserProfile) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	u.CreatedAt = time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *InMemoryRepository) GetUserByID(ctx context.Context, id string) (*UserProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *InMemoryRepository) ListUsers(ctx context.Context) ([]*UserProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*UserProfile
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

func (r *InMemoryRepository) AddToSegment(ctx context.Context, userID string, segmentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[userID]
	if !ok {
		return nil
	}
	u.SegmentIDs = append(u.SegmentIDs, segmentID)
	return nil
}
