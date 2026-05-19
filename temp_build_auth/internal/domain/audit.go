package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type AuditAction string

const (
	AuditLogin          AuditAction = "login"
	AuditLogout         AuditAction = "logout"
	AuditLoginFailed    AuditAction = "login_failed"
	AuditTokenRefresh   AuditAction = "token_refresh"
	AuditTokenRevoke    AuditAction = "token_revoke"
	AuditRegister       AuditAction = "register"
	AuditPasswordReset  AuditAction = "password_reset"
	AuditPasswordChange AuditAction = "password_change"
	AuditEmailVerify    AuditAction = "email_verify"
	AuditRoleChange     AuditAction = "role_change"
	AuditSessionRevoke  AuditAction = "session_revoke"
	AuditAccountLock    AuditAction = "account_lock"
	AuditMFAEnable      AuditAction = "mfa_enable"
	AuditMFADisable     AuditAction = "mfa_disable"
	AuditProfileUpdate  AuditAction = "profile_update"
	AuditSuspicious     AuditAction = "suspicious_activity"
)

type AuditLog struct {
	ID        string      `db:"id" json:"id"`
	TraceID   string      `db:"trace_id" json:"trace_id"`
	UserID    string      `db:"user_id" json:"user_id,omitempty"`
	Action    AuditAction `db:"action" json:"action"`
	IP        string      `db:"ip" json:"ip"`
	DeviceID  string      `db:"device_id" json:"device_id,omitempty"`
	UserAgent string      `db:"user_agent" json:"user_agent,omitempty"`
	Resource  string      `db:"resource" json:"resource,omitempty"`
	Status    string      `db:"status" json:"status"`
	Detail    string      `db:"detail" json:"detail,omitempty"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
}

func NewAuditLog(traceID, userID string, action AuditAction, ip, deviceID, userAgent string) *AuditLog {
	return &AuditLog{
		ID:        generateID(),
		TraceID:   traceID,
		UserID:    userID,
		Action:    action,
		IP:        ip,
		DeviceID:  deviceID,
		UserAgent: userAgent,
		Status:    "success",
		CreatedAt: time.Now(),
	}
}

func generateID() string {
	return "aud_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(n int) string {
	b := make([]byte, n/2+1)
	if _, err := rand.Read(b); err != nil {
		// fallback: use low-entropy time-based approach
		for i := range b {
			b[i] = byte(time.Now().UnixNano() >> uint(i*8))
		}
	}
	return hex.EncodeToString(b)[:n]
}
