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
	ID             string        `db:"id" json:"id"`
	UserID         string        `db:"user_id" json:"user_id"`
	RefreshToken   string        `db:"-" json:"-"`
	RefreshTokenID string        `db:"refresh_token_id" json:"-"`
	DeviceID       string        `db:"device_id" json:"device_id,omitempty"`
	DeviceName     string        `db:"device_name" json:"device_name,omitempty"`
	DeviceType     string        `db:"device_type" json:"device_type,omitempty"`
	Platform       string        `db:"platform" json:"platform,omitempty"`
	IP             string        `db:"ip" json:"ip"`
	UserAgent      string        `db:"user_agent" json:"user_agent,omitempty"`
	Location       string        `db:"location" json:"location,omitempty"`
	Status         SessionStatus `db:"status" json:"status"`
	LastActiveAt   time.Time     `db:"last_active_at" json:"last_active_at"`
	ExpiresAt      time.Time     `db:"expires_at" json:"expires_at"`
	CreatedAt      time.Time     `db:"created_at" json:"created_at"`
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
