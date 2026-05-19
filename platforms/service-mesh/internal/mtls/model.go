package mtls

import "time"

type Certificate struct {
	ID           string    `json:"id"`
	ServiceName  string    `json:"service_name"`
	CommonName   string    `json:"common_name"`
	Organization string    `json:"organization"`
	ValidityDays int       `json:"validity_days"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Serial       string    `json:"serial"`
	Fingerprint  string    `json:"fingerprint"`
	IsRevoked    bool      `json:"is_revoked"`
	IsCA         bool      `json:"is_ca"`
}
