package domain

import (
	"context"
	"database/sql"
)

type ServiceHealthRepository interface {
	FindByID(ctx context.Context, id string) (*ServiceHealth, error)
	FindByServiceName(ctx context.Context, name string) (*ServiceHealth, error)
	FindAll(ctx context.Context) ([]*ServiceHealth, error)
	FindByStatus(ctx context.Context, status string) ([]*ServiceHealth, error)
	Create(ctx context.Context, health *ServiceHealth) error
	Update(ctx context.Context, health *ServiceHealth) error
	Delete(ctx context.Context, id string) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type DeploymentRepository interface {
	FindByID(ctx context.Context, id string) (*Deployment, error)
	FindByServiceName(ctx context.Context, serviceName string, limit int) ([]*Deployment, error)
	FindActive(ctx context.Context) ([]*Deployment, error)
	FindRecent(ctx context.Context, limit int) ([]*Deployment, error)
	Create(ctx context.Context, deployment *Deployment) error
	Update(ctx context.Context, deployment *Deployment) error
	Delete(ctx context.Context, id string) error
	CountByStatus(ctx context.Context, status string) (int, error)
}

type IncidentRepository interface {
	FindByID(ctx context.Context, id string) (*Incident, error)
	FindActive(ctx context.Context) ([]*Incident, error)
	FindBySeverity(ctx context.Context, severity string) ([]*Incident, error)
	FindByStatus(ctx context.Context, status string) ([]*Incident, error)
	FindByServiceName(ctx context.Context, serviceName string) ([]*Incident, error)
	FindRecent(ctx context.Context, limit int) ([]*Incident, error)
	Create(ctx context.Context, incident *Incident) error
	Update(ctx context.Context, incident *Incident) error
	Delete(ctx context.Context, id string) error
	CountActive(ctx context.Context) (int, error)
	CountBySeverity(ctx context.Context, severity string) (int, error)
	CountActiveBySeverity(ctx context.Context) (map[string]int, error)
}

type AlertRuleRepository interface {
	FindByID(ctx context.Context, id string) (*AlertRule, error)
	FindAll(ctx context.Context) ([]*AlertRule, error)
	FindByServiceName(ctx context.Context, serviceName string) ([]*AlertRule, error)
	FindEnabled(ctx context.Context) ([]*AlertRule, error)
	Create(ctx context.Context, rule *AlertRule) error
	Update(ctx context.Context, rule *AlertRule) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type AuditLogRepository interface {
	FindByID(ctx context.Context, id string) (*AuditLog, error)
	FindByActor(ctx context.Context, actor string, limit int) ([]*AuditLog, error)
	FindByResource(ctx context.Context, resource, resourceID string, limit int) ([]*AuditLog, error)
	FindRecent(ctx context.Context, limit int) ([]*AuditLog, error)
	Create(ctx context.Context, log *AuditLog) error
}

type ServiceDependencyRepository interface {
	FindByService(ctx context.Context, serviceName string) ([]*ServiceDependency, error)
	FindAll(ctx context.Context) ([]*ServiceDependency, error)
	FindDependents(ctx context.Context, serviceName string) ([]*ServiceDependency, error)
	Create(ctx context.Context, dep *ServiceDependency) error
	Delete(ctx context.Context, id string) error
	DeleteByService(ctx context.Context, serviceName string) error
}

type CapacityMetricRepository interface {
	FindByService(ctx context.Context, serviceName string) ([]*CapacityMetric, error)
	FindByResource(ctx context.Context, serviceName, resourceType string, limit int) ([]*CapacityMetric, error)
	FindLatest(ctx context.Context) ([]*CapacityMetric, error)
	Create(ctx context.Context, metric *CapacityMetric) error
}
