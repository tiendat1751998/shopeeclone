package domain

import (
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionActive    SessionStatus = "active"
	SessionExpired   SessionStatus = "expired"
	SessionRevoked   SessionStatus = "revoked"
	SessionReplaced  SessionStatus = "replaced"
)

type Session struct {
	ID             string        `json:"id"`
	UserID         string        `json:"user_id"`
	RefreshToken   string        `json:"-"`
	RefreshTokenID string        `json:"-"` // jti
	DeviceID       string        `json:"device_id,omitempty"`
	DeviceName     string        `json:"device_name,omitempty"`
	DeviceType     string        `json:"device_type,omitempty"`
	Platform       string        `json:"platform,omitempty"`
	IP             string        `json:"ip"`
	UserAgent      string        `json:"user_agent,omitempty"`
	Location       string        `json:"location,omitempty"`
	Status         SessionStatus `json:"status"`
	LastActiveAt   time.Time     `json:"last_active_at"`
	ExpiresAt      time.Time     `json:"expires_at"`
	CreatedAt      time.Time     `json:"created_at"`
}

func NewSession(userID, ip, userAgent string, refreshTokenID string, ttl time.Duration) *Session {
	now := time.Now()
	return &Session{
		ID:             uuid.New().String(),
		UserID:         userID,
		RefreshTokenID: refreshTokenID,
		IP:             ip,
		UserAgent:      userAgent,
		Status:         SessionActive,
		LastActiveAt:   now,
		ExpiresAt:      now.Add(ttl),
		CreatedAt:      now,
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Expire() {
	s.Status = SessionExpired
}

func (s *Session) Revoke() {
	s.Status = SessionRevoked
}

func (s *Session) Touch() {
	s.LastActiveAt = time.Now()
}
