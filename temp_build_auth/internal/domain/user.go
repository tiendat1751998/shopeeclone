package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusLocked   UserStatus = "locked"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending  UserStatus = "pending"
)

type User struct {
	ID             string     `db:"id" json:"id"`
	Email          string     `db:"email" json:"email"`
	Phone          string     `db:"phone" json:"phone,omitempty"`
	Username       string     `db:"username" json:"username"`
	PasswordHash   string     `db:"password_hash" json:"-"`
	DisplayName    string     `db:"display_name" json:"display_name"`
	Status         UserStatus `db:"status" json:"status"`
	EmailVerified  bool       `db:"email_verified" json:"email_verified"`
	PhoneVerified  bool       `db:"phone_verified" json:"phone_verified"`
	MFAEnabled     bool       `db:"mfa_enabled" json:"mfa_enabled"`
	TwoFASecret    string     `db:"twofa_secret" json:"-"`
	LastLoginAt    *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	LastLoginIP    string     `db:"last_login_ip" json:"-"`
	FailedAttempts int        `db:"failed_attempts" json:"-"`
	LockedUntil    *time.Time `db:"locked_until" json:"-"`
	Metadata       string     `db:"metadata" json:"-"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

func NewUser(email, username, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:            uuid.New().String(),
		Email:         email,
		Username:      username,
		PasswordHash:  passwordHash,
		Status:        UserStatusPending,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

func (u *User) CanLogin() error {
	switch u.Status {
	case UserStatusInactive:
		return ErrAccountInactive
	case UserStatusLocked:
		return ErrAccountLocked
	case UserStatusSuspended:
		return ErrAccountSuspended
	case UserStatusPending:
		if !u.EmailVerified {
			return ErrEmailNotVerified
		}
	}
	if u.IsLocked() {
		return ErrAccountLocked
	}
	return nil
}

func (u *User) RecordFailedAttempt() {
	u.FailedAttempts++
	u.UpdatedAt = time.Now()
}

func (u *User) RecordSuccessfulLogin(ip string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastLoginIP = ip
	u.FailedAttempts = 0
	u.LockedUntil = nil
	u.UpdatedAt = now
}

func (u *User) Lock(duration time.Duration) {
	until := time.Now().Add(duration)
	u.LockedUntil = &until
	u.Status = UserStatusLocked
	u.UpdatedAt = time.Now()
}
