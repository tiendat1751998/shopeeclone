package discovery

import "time"

type Status string

const (
	StatusUp      Status = "up"
	StatusDown    Status = "down"
	StatusDraining Status = "draining"
)

type ServiceInstance struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Address         string            `json:"address"`
	Port            int               `json:"port"`
	Region          string            `json:"region"`
	Zone            string            `json:"zone"`
	Status          Status            `json:"status"`
	Metadata        map[string]string `json:"metadata"`
	HealthCheckPath string            `json:"health_check_path"`
	LastHeartbeat   time.Time         `json:"last_heartbeat"`
}
