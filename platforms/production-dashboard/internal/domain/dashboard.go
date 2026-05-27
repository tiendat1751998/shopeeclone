package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// === Service Health ===

type ServiceHealth struct {
	ID            string    `db:"id" json:"id"`
	ServiceName   string    `db:"service_name" json:"service_name"`
	Status        string    `db:"status" json:"status"`
	HealthURL     string    `db:"health_url" json:"health_url"`
	ResponseTime  int       `db:"response_time_ms" json:"response_time_ms"`
	LastCheckedAt time.Time `db:"last_checked_at" json:"last_checked_at"`
	LastError     string    `db:"last_error" json:"last_error,omitempty"`
	Version       string    `db:"version" json:"version,omitempty"`
	Environment   string    `db:"environment" json:"environment"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

const (
	ServiceStatusHealthy   = "healthy"
	ServiceStatusDegraded  = "degraded"
	ServiceStatusUnhealthy = "unhealthy"
	ServiceStatusUnknown   = "unknown"
)

func NewServiceHealth(name, healthURL, environment string) *ServiceHealth {
	now := time.Now()
	return &ServiceHealth{
		ID:          uuid.New().String(),
		ServiceName: name,
		Status:      ServiceStatusUnknown,
		HealthURL:   healthURL,
		Environment: environment,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (s *ServiceHealth) MarkHealthy(responseTime int, version string) {
	s.Status = ServiceStatusHealthy
	s.ResponseTime = responseTime
	s.LastError = ""
	s.Version = version
	s.LastCheckedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *ServiceHealth) MarkDegraded(responseTime int, errMsg string) {
	s.Status = ServiceStatusDegraded
	s.ResponseTime = responseTime
	s.LastError = errMsg
	s.LastCheckedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *ServiceHealth) MarkUnhealthy(errMsg string) {
	s.Status = ServiceStatusUnhealthy
	s.LastError = errMsg
	s.ResponseTime = 0
	s.LastCheckedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *ServiceHealth) IsHealthy() bool {
	return s.Status == ServiceStatusHealthy
}

// === Deployment ===

type Deployment struct {
	ID            string    `db:"id" json:"id"`
	ServiceName   string    `db:"service_name" json:"service_name"`
	Version       string    `db:"version" json:"version"`
	Environment   string    `db:"environment" environment"`
	Status        string    `db:"status" json:"status"`
	DeployedBy    string    `db:"deployed_by" json:"deployed_by"`
	Image         string    `db:"image" json:"image"`
	Replicas      int       `db:"replicas" json:"replicas"`
	ReadyReplicas int       `db:"ready_replicas" json:"ready_replicas"`
	Strategy      string    `db:"strategy" json:"strategy"`
	StartedAt     time.Time `db:"started_at" json:"started_at"`
	FinishedAt    time.Time `db:"finished_at" json:"finished_at,omitempty"`
	RollbackOf    string    `db:"rollback_of" json:"rollback_of,omitempty"`
	Notes         string    `db:"notes" json:"notes,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

const (
	DeploymentStatusPending    = "pending"
	DeploymentStatusInProgress = "in_progress"
	DeploymentStatusSucceeded  = "succeeded"
	DeploymentStatusFailed     = "failed"
	DeploymentStatusRolledBack = "rolled_back"
)

func NewDeployment(serviceName, version, environment, deployedBy, image string, replicas int, strategy string) *Deployment {
	now := time.Now()
	return &Deployment{
		ID:          uuid.New().String(),
		ServiceName: serviceName,
		Version:     version,
		Environment: environment,
		Status:      DeploymentStatusPending,
		DeployedBy:  deployedBy,
		Image:       image,
		Replicas:    replicas,
		Strategy:    strategy,
		StartedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (d *Deployment) MarkInProgress() {
	d.Status = DeploymentStatusInProgress
	d.UpdatedAt = time.Now()
}

func (d *Deployment) MarkSucceeded(readyReplicas int) {
	d.Status = DeploymentStatusSucceeded
	d.ReadyReplicas = readyReplicas
	d.FinishedAt = time.Now()
	d.UpdatedAt = time.Now()
}

func (d *Deployment) MarkFailed() {
	d.Status = DeploymentStatusFailed
	d.FinishedAt = time.Now()
	d.UpdatedAt = time.Now()
}

func (d *Deployment) MarkRolledBack() {
	d.Status = DeploymentStatusRolledBack
	d.FinishedAt = time.Now()
	d.UpdatedAt = time.Now()
}

func (d *Deployment) IsActive() bool {
	return d.Status == DeploymentStatusPending || d.Status == DeploymentStatusInProgress
}

// === Incident ===

type Incident struct {
	ID             string    `db:"id" json:"id"`
	Title          string    `db:"title" json:"title"`
	Description    string    `db:"description" json:"description"`
	Severity       string    `db:"severity" json:"severity"`
	Status         string    `db:"status" json:"status"`
	ServiceNames   string    `db:"service_names" json:"service_names"`
	DetectedAt     time.Time `db:"detected_at" json:"detected_at"`
	AcknowledgedAt time.Time `db:"acknowledged_at" json:"acknowledged_at,omitempty"`
	ResolvedAt     time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
	RootCause      string    `db:"root_cause" json:"root_cause,omitempty"`
	Resolution     string    `db:"resolution" json:"resolution,omitempty"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	AssignedTo     string    `db:"assigned_to" json:"assigned_to,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

const (
	IncidentSeverityCritical = "critical"
	IncidentSeverityHigh     = "high"
	IncidentSeverityMedium   = "medium"
	IncidentSeverityLow      = "low"

	IncidentStatusOpen        = "open"
	IncidentStatusAcknowledged = "acknowledged"
	IncidentStatusInvestigating = "investigating"
	IncidentStatusMitigating = "mitigating"
	IncidentStatusResolved   = "resolved"
	IncidentStatusClosed     = "closed"
)

func NewIncident(title, description, severity, serviceNames, createdBy string) *Incident {
	now := time.Now()
	return &Incident{
		ID:           uuid.New().String(),
		Title:        title,
		Description:  description,
		Severity:     severity,
		Status:       IncidentStatusOpen,
		ServiceNames: serviceNames,
		DetectedAt:   now,
		CreatedBy:    createdBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (i *Incident) Acknowledge(assignedTo string) error {
	if i.Status != IncidentStatusOpen {
		return fmt.Errorf("%w: incident must be open to acknowledge", ErrInvalidIncidentState)
	}
	i.Status = IncidentStatusAcknowledged
	i.AssignedTo = assignedTo
	i.AcknowledgedAt = time.Now()
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Incident) StartInvestigation() error {
	if i.Status != IncidentStatusAcknowledged && i.Status != IncidentStatusOpen {
		return fmt.Errorf("%w: incident must be open or acknowledged to investigate", ErrInvalidIncidentState)
	}
	i.Status = IncidentStatusInvestigating
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Incident) Mitigate() error {
	if i.Status != IncidentStatusInvestigating {
		return fmt.Errorf("%w: incident must be investigating to mitigate", ErrInvalidIncidentState)
	}
	i.Status = IncidentStatusMitigating
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Incident) Resolve(rootCause, resolution string) error {
	if i.Status == IncidentStatusResolved || i.Status == IncidentStatusClosed {
		return fmt.Errorf("%w: incident already resolved", ErrInvalidIncidentState)
	}
	i.Status = IncidentStatusResolved
	i.RootCause = rootCause
	i.Resolution = resolution
	i.ResolvedAt = time.Now()
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Incident) Close() error {
	if i.Status != IncidentStatusResolved {
		return fmt.Errorf("%w: incident must be resolved before closing", ErrInvalidIncidentState)
	}
	i.Status = IncidentStatusClosed
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Incident) Duration() time.Duration {
	end := i.ResolvedAt
	if end.IsZero() {
		end = time.Now()
	}
	return end.Sub(i.DetectedAt)
}

// === Alert Rule ===

type AlertRule struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	ServiceName string    `db:"service_name" json:"service_name"`
	MetricName  string    `db:"metric_name" json:"metric_name"`
	Condition   string    `db:"condition" json:"condition"`
	Threshold   float64   `db:"threshold" json:"threshold"`
	Duration    string    `db:"duration" json:"duration"`
	Severity    string    `db:"severity" json:"severity"`
	Enabled     bool      `db:"enabled" json:"enabled"`
	NotifyChannels string `db:"notify_channels" json:"notify_channels"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

const (
	AlertConditionGT  = "gt"
	AlertConditionLT  = "lt"
	AlertConditionGTE = "gte"
	AlertConditionLTE = "lte"
	AlertConditionEQ  = "eq"
)

func NewAlertRule(name, description, serviceName, metricName, condition string, threshold float64, duration, severity, notifyChannels, createdBy string) *AlertRule {
	now := time.Now()
	return &AlertRule{
		ID:             uuid.New().String(),
		Name:           name,
		Description:    description,
		ServiceName:    serviceName,
		MetricName:     metricName,
		Condition:      condition,
		Threshold:      threshold,
		Duration:       duration,
		Severity:       severity,
		Enabled:        true,
		NotifyChannels: notifyChannels,
		CreatedBy:      createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (a *AlertRule) Evaluate(value float64) bool {
	if !a.Enabled {
		return false
	}
	switch a.Condition {
	case AlertConditionGT:
		return value > a.Threshold
	case AlertConditionLT:
		return value < a.Threshold
	case AlertConditionGTE:
		return value >= a.Threshold
	case AlertConditionLTE:
		return value <= a.Threshold
	case AlertConditionEQ:
		return value == a.Threshold
	default:
		return false
	}
}

// === Audit Log ===

type AuditLog struct {
	ID         string `db:"id" json:"id"`
	Actor      string `db:"actor" json:"actor"`
	Action     string `db:"action" json:"action"`
	Resource   string `db:"resource" json:"resource"`
	ResourceID string `db:"resource_id" json:"resource_id,omitempty"`
	Details    string `db:"details" json:"details,omitempty"`
	IPAddress  string `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent  string `db:"user_agent" json:"user_agent,omitempty"`
	CreatedAt  string `db:"created_at" json:"created_at"`
}

const (
	ActionCreate   = "create"
	ActionUpdate   = "update"
	ActionDelete   = "delete"
	ActionDeploy   = "deploy"
	ActionRollback = "rollback"
	ActionAcknowledge = "acknowledge"
	ActionResolve  = "resolve"
	ActionClose    = "close"
	ActionScale    = "scale"
	ActionRestart  = "restart"
)

func NewAuditLog(actor, action, resource, resourceID, details, ipAddress, userAgent string) *AuditLog {
	return &AuditLog{
		ID:         uuid.New().String(),
		Actor:      actor,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
}

// === System Topology ===

type ServiceDependency struct {
	ID             string `db:"id" json:"id"`
	ServiceName    string `db:"service_name" json:"service_name"`
	DependsOn      string `db:"depends_on" json:"depends_on"`
	DependencyType string `db:"dependency_type" json:"dependency_type"`
	Critical       bool   `db:"critical" json:"critical"`
	CreatedAt      string `db:"created_at" json:"created_at"`
}

const (
	DependencyTypeSync   = "sync"
	DependencyTypeAsync  = "async"
	DependencyTypeData   = "data"
)

func NewServiceDependency(serviceName, dependsOn, depType string, critical bool) *ServiceDependency {
	return &ServiceDependency{
		ID:             uuid.New().String(),
		ServiceName:    serviceName,
		DependsOn:      dependsOn,
		DependencyType: depType,
		Critical:       critical,
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
	}
}

// === Capacity Metric ===

type CapacityMetric struct {
	ID           string  `db:"id" json:"id"`
	ServiceName  string  `db:"service_name" json:"service_name"`
	ResourceType string  `db:"resource_type" json:"resource_type"`
	CurrentValue float64 `db:"current_value" json:"current_value"`
	MaxValue     float64 `db:"max_value" json:"max_value"`
	Unit         string  `db:"unit" json:"unit"`
	RecordedAt   string  `db:"recorded_at" json:"recorded_at"`
}

const (
	ResourceCPU     = "cpu"
	ResourceMemory  = "memory"
	ResourceDisk    = "disk"
	ResourceNetwork = "network"
	ResourcePods    = "pods"
)

func NewCapacityMetric(serviceName, resourceType string, current, max float64, unit string) *CapacityMetric {
	return &CapacityMetric{
		ID:           uuid.New().String(),
		ServiceName:  serviceName,
		ResourceType: resourceType,
		CurrentValue: current,
		MaxValue:     max,
		Unit:         unit,
		RecordedAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

func (c *CapacityMetric) UtilizationPercent() float64 {
	if c.MaxValue == 0 {
		return 0
	}
	return (c.CurrentValue / c.MaxValue) * 100
}

// === Dashboard Summary ===

type DashboardSummary struct {
	TotalServices      int                    `json:"total_services"`
	HealthyServices    int                    `json:"healthy_services"`
	DegradedServices   int                    `json:"degraded_services"`
	UnhealthyServices  int                    `json:"unhealthy_services"`
	ActiveIncidents    int                    `json:"active_incidents"`
	CriticalIncidents  int                    `json:"critical_incidents"`
	ActiveDeployments  int                    `json:"active_deployments"`
	RecentDeployments  []*Deployment           `json:"recent_deployments"`
	ServiceHealthList  []*ServiceHealth        `json:"service_health_list"`
	OpenIncidents      []*Incident             `json:"open_incidents"`
	CapacityOverview   []*CapacityMetric       `json:"capacity_overview"`
}

// === Domain Errors ===

var (
	ErrServiceNotFound    = ErrDashboard("service_not_found")
	ErrDeploymentNotFound = ErrDashboard("deployment_not_found")
	ErrIncidentNotFound   = ErrDashboard("incident_not_found")
	ErrAlertRuleNotFound  = ErrDashboard("alert_rule_not_found")
	ErrInvalidIncidentState = ErrDashboard("invalid_incident_state")
	ErrInvalidDeploymentState = ErrDashboard("invalid_deployment_state")
	ErrDuplicateService   = ErrDashboard("duplicate_service")
	ErrDuplicateAlertRule = ErrDashboard("duplicate_alert_rule")
)

type ErrDashboard string

func (e ErrDashboard) Error() string {
	return "dashboard: " + string(e)
}
