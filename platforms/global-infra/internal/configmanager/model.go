package configmanager

import "time"

type Environment string

const (
	EnvDev     Environment = "dev"
	EnvStaging Environment = "staging"
	EnvProd    Environment = "prod"
)

type ConfigEntry struct {
	Key          string      `json:"key"`
	Value        string      `json:"value"`
	Environment  Environment `json:"environment"`
	ServiceName  string      `json:"service_name"`
	Version      int         `json:"version"`
	IsEncrypted  bool        `json:"is_encrypted"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type ConfigChangeEvent struct {
	Key         string      `json:"key"`
	Environment Environment `json:"environment"`
	ServiceName string      `json:"service_name"`
	Version     int         `json:"version"`
	Action      string      `json:"action"`
	Timestamp   time.Time   `json:"timestamp"`
}
