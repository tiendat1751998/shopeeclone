package preferences

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	Get(ctx context.Context, userID string) (*UserPreference, error)
	Upsert(ctx context.Context, pref *UserPreference) error
	AddSuppression(ctx context.Context, entry *SuppressionEntry) error
	IsSuppressed(ctx context.Context, userID, email, phone string) (bool, error)
}

type InMemoryRepository struct {
	mu           sync.RWMutex
	preferences  map[string]*UserPreference
	suppressions []*SuppressionEntry
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		preferences:  make(map[string]*UserPreference),
		suppressions: make([]*SuppressionEntry, 0),
	}
}

func (r *InMemoryRepository) Get(ctx context.Context, userID string) (*UserPreference, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pref, ok := r.preferences[userID]
	if !ok {
		return &UserPreference{
			UserID: userID,
			ChannelOptIn: ChannelOptIn{
				Push:  true,
				Email: true,
				SMS:   true,
				InApp: true,
			},
			Categories:    CategoryPreferences{},
			PushEnabled:   true,
			EmailDigest:   false,
			SMSPromotions: true,
		}, nil
	}
	return pref, nil
}

func (r *InMemoryRepository) Upsert(ctx context.Context, pref *UserPreference) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pref.UpdatedAt = time.Now()
	r.preferences[pref.UserID] = pref
	return nil
}

func (r *InMemoryRepository) AddSuppression(ctx context.Context, entry *SuppressionEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry.CreatedAt = time.Now()
	r.suppressions = append(r.suppressions, entry)
	return nil
}

func (r *InMemoryRepository) IsSuppressed(ctx context.Context, userID, email, phone string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, entry := range r.suppressions {
		if entry.UserID == userID {
			return true, nil
		}
		if email != "" && entry.Email == email {
			return true, nil
		}
		if phone != "" && entry.Phone == phone {
			return true, nil
		}
	}
	return false, nil
}
