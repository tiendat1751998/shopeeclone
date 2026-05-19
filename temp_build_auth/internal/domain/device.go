package domain

import "time"

type DeviceInfo struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	Platform   string `json:"platform"`
	OS         string `json:"os"`
	Browser    string `json:"browser"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Fingerprint string `json:"fingerprint,omitempty"`
	IsTrusted  bool   `json:"is_trusted"`
	LastUsedAt time.Time `json:"last_used_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type TokenPayload struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Roles     []Role   `json:"roles"`
	SessionID string   `json:"session_id"`
	DeviceID  string   `json:"device_id"`
	TokenID   string   `json:"jti"`
	Type      string   `json:"type"`
}

type LoginAttempt struct {
	Email     string    `json:"email"`
	IP        string    `json:"ip"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"user_agent"`
	Reason    string    `json:"reason,omitempty"`
}

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Username        string `json:"username"`
	DisplayName     string `json:"display_name"`
	Phone           string `json:"phone,omitempty"`
	DeviceID        string `json:"device_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	SessionID    string `json:"session_id"`
}
