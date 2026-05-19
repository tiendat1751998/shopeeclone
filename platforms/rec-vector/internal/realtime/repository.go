package realtime

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreSession(ctx context.Context, session *UserSession) error
	GetSession(ctx context.Context, sessionID string) (*UserSession, error)
	UpdateSession(ctx context.Context, session *UserSession) error
	StoreArmStats(ctx context.Context, stats map[string]*ArmStat) error
	GetArmStats(ctx context.Context) (map[string]*ArmStat, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	sessions map[string]*UserSession
	armStats map[string]*ArmStat
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		sessions: make(map[string]*UserSession),
		armStats: make(map[string]*ArmStat),
	}
}

func (r *InMemoryRepository) StoreSession(ctx context.Context, session *UserSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if session.StartedAt.IsZero() {
		session.StartedAt = time.Now()
	}
	session.LastActiveAt = time.Now()
	r.sessions[session.SessionID] = session
	return nil
}

func (r *InMemoryRepository) GetSession(ctx context.Context, sessionID string) (*UserSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return s, nil
}

func (r *InMemoryRepository) UpdateSession(ctx context.Context, session *UserSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	session.LastActiveAt = time.Now()
	r.sessions[session.SessionID] = session
	return nil
}

func (r *InMemoryRepository) StoreArmStats(ctx context.Context, stats map[string]*ArmStat) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, v := range stats {
		r.armStats[k] = v
	}
	return nil
}

func (r *InMemoryRepository) GetArmStats(ctx context.Context) (map[string]*ArmStat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]*ArmStat)
	for k, v := range r.armStats {
		result[k] = &ArmStat{
			ArmID:      v.ArmID,
			Plays:      v.Plays,
			Rewards:    v.Rewards,
			MeanReward: v.MeanReward,
		}
	}
	return result, nil
}
