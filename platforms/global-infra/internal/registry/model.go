package registry

import "time"

type Status string

const (
	StatusUp      Status = "up"
	StatusDown    Status = "down"
	StatusDraining Status = "draining"
)

type ServiceInstance struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Address        string            `json:"address"`
	Port           int               `json:"port"`
	Region         string            `json:"region"`
	Status         Status            `json:"status"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	HealthEndpoint string            `json:"health_endpoint"`
	LastHeartbeat  time.Time         `json:"last_heartbeat"`
	RegisteredAt   time.Time         `json:"registered_at"`
}
