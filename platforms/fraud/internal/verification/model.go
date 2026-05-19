package verification

import "time"

type VerificationMethod string

const (
	MethodSMS   VerificationMethod = "sms"
	MethodEmail VerificationMethod = "email"
	MethodKYC   VerificationMethod = "kyc"
)

type VerificationStatus string

const (
	StatusPending  VerificationStatus = "pending"
	StatusVerified VerificationStatus = "verified"
	StatusExpired  VerificationStatus = "expired"
	StatusFailed   VerificationStatus = "failed"
)

type VerificationRequest struct {
	ID            string             `json:"id"`
	UserID        string             `json:"user_id"`
	Method        VerificationMethod `json:"method"`
	Target        string             `json:"target"`
	Code          string             `json:"-"`
	Status        VerificationStatus `json:"status"`
	Attempts      int                `json:"attempts"`
	MaxAttempts   int                `json:"max_attempts"`
	CreatedAt     time.Time          `json:"created_at"`
	ExpiresAt     time.Time          `json:"expires_at"`
	VerifiedAt    *time.Time         `json:"verified_at,omitempty"`
}

type KYCStatus struct {
	UserID        string `json:"user_id"`
	IsVerified    bool   `json:"is_verified"`
	Level         string `json:"level"`
	DocumentType  string `json:"document_type"`
	SubmittedAt   *time.Time `json:"submitted_at,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	RejectedReason string `json:"rejected_reason,omitempty"`
}
