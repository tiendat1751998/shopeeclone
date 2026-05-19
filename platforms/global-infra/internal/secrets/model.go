package secrets

import "time"

type Secret struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Value           string    `json:"-"`
	ServiceName     string    `json:"service_name"`
	RotationPeriod  int       `json:"rotation_period"`
	LastRotated     time.Time `json:"last_rotated"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
