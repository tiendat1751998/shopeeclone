package session

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreSession(ctx context.Context, session *Session) error
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	ListSessions(ctx context.Context, filter *SessionFilter) ([]*Session, int64, error)
	UpdateSession(ctx context.Context, session *Session) error
	GetActiveSession(ctx context.Context, userID string) (*Session, error)
	GetSessionMetrics(ctx context.Context, filter *SessionFilter) (*SessionMetrics, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		sessions: make(map[string]*Session),
	}
}

func (r *InMemoryRepository) StoreSession(ctx context.Context, session *Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	r.sessions[session.SessionID] = session
	return nil
}

func (r *InMemoryRepository) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *InMemoryRepository) ListSessions(ctx context.Context, filter *SessionFilter) ([]*Session, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*Session
	for _, s := range r.sessions {
		if filter.UserID != "" && s.UserID != filter.UserID {
			continue
		}
		if filter.Source != "" && s.Source != filter.Source {
			continue
		}
		if filter.Device != "" && s.Device != filter.Device {
			continue
		}
		if filter.Country != "" && s.Country != filter.Country {
			continue
		}
		if filter.StartTime != nil && s.StartTime.Before(*filter.StartTime) {
			continue
		}
		if filter.EndTime != nil && s.StartTime.After(*filter.EndTime) {
			continue
		}
		if filter.MinDuration > 0 && int64(s.Duration.Seconds()) < filter.MinDuration {
			continue
		}
		if filter.HasConversion != nil && s.HasConversion != *filter.HasConversion {
			continue
		}
		filtered = append(filtered, s)
	}
	total := int64(len(filtered))
	offset := filter.Offset
	limit := filter.Limit
	if offset >= len(filtered) {
		return []*Session{}, total, nil
	}
	end := offset + limit
	if end > len(filtered) || limit == 0 {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (r *InMemoryRepository) UpdateSession(ctx context.Context, session *Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.SessionID] = session
	return nil
}

func (r *InMemoryRepository) GetActiveSession(ctx context.Context, userID string) (*Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.sessions {
		if s.UserID == userID && s.IsActive {
			return s, nil
		}
	}
	return nil, nil
}

func (r *InMemoryRepository) GetSessionMetrics(ctx context.Context, filter *SessionFilter) (*SessionMetrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metrics := &SessionMetrics{}
	var totalDuration float64
	var totalPageviews int64
	var bounceCount int64
	var conversionCount int64

	for _, s := range r.sessions {
		if filter.Source != "" && s.Source != filter.Source {
			continue
		}
		if filter.Device != "" && s.Device != filter.Device {
			continue
		}
		if filter.StartTime != nil && s.StartTime.Before(*filter.StartTime) {
			continue
		}
		if filter.EndTime != nil && s.StartTime.After(*filter.EndTime) {
			continue
		}

		metrics.TotalSessions++
		if s.IsActive {
			metrics.ActiveSessions++
		}
		totalDuration += s.Duration.Seconds()
		totalPageviews += s.Pageviews
		if s.Pageviews <= 1 {
			bounceCount++
		}
		if s.HasConversion {
			conversionCount++
		}
		metrics.TotalRevenue += s.Revenue
	}

	if metrics.TotalSessions > 0 {
		metrics.AvgDuration = totalDuration / float64(metrics.TotalSessions)
		metrics.AvgPageviews = float64(totalPageviews) / float64(metrics.TotalSessions)
		metrics.BounceRate = float64(bounceCount) / float64(metrics.TotalSessions) * 100
		metrics.ConversionRate = float64(conversionCount) / float64(metrics.TotalSessions) * 100
		metrics.AvgSessionRevenue = metrics.TotalRevenue / float64(metrics.TotalSessions)
	}

	return metrics, nil
}
