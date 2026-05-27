package domain

import "time"

type DashboardEvent struct {
	ID            string      `json:"id"`
	EventType     string      `json:"event_type"`
	AggregateType string      `json:"aggregate_type"`
	AggregateID   string      `json:"aggregate_id"`
	Payload       interface{} `json:"payload"`
	CreatedAt     time.Time   `json:"created_at"`
}

const (
	EventIncidentCreated         = "dashboard.incident_created"
	EventIncidentAcknowledged   = "dashboard.incident_acknowledged"
	EventIncidentResolved       = "dashboard.incident_resolved"
	EventDeploymentStarted      = "dashboard.deployment_started"
	EventDeploymentSucceeded    = "dashboard.deployment_succeeded"
	EventDeploymentFailed       = "dashboard.deployment_failed"
	EventDeploymentRolledBack   = "dashboard.deployment_rolled_back"
	EventServiceStatusChanged   = "dashboard.service_status_changed"
	EventAlertFired             = "dashboard.alert_fired"
	EventAlertResolved          = "dashboard.alert_resolved"
)

type IncidentCreatedPayload struct {
	IncidentID string `json:"incident_id"`
	Title      string `json:"title"`
	Severity   string `json:"severity"`
	Services   string `json:"services"`
}

type IncidentResolvedPayload struct {
	IncidentID string `json:"incident_id"`
	RootCause  string `json:"root_cause"`
	Duration   string `json:"duration"`
}

type DeploymentEventPayload struct {
	DeploymentID string `json:"deployment_id"`
	ServiceName  string `json:"service_name"`
	Version      string `json:"version"`
	Environment  string `json:"environment"`
}

type ServiceStatusChangedPayload struct {
	ServiceName    string `json:"service_name"`
	PreviousStatus string `json:"previous_status"`
	CurrentStatus  string `json:"current_status"`
}

type AlertFiredPayload struct {
	AlertRuleID string  `json:"alert_rule_id"`
	ServiceName string  `json:"service_name"`
	MetricName  string  `json:"metric_name"`
	Value       float64 `json:"value"`
	Threshold   float64 `json:"threshold"`
}
