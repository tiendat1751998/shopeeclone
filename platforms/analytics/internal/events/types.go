package events

import "time"

type PlatformEventType string

const (
	EventReportGenerated   PlatformEventType = "analytics.report.generated"
	EventDashboardShared   PlatformEventType = "analytics.dashboard.shared"
	EventAnomalyDetected   PlatformEventType = "analytics.anomaly.detected"
	EventScheduledReport   PlatformEventType = "analytics.scheduled.report"
	EventAlertTriggered    PlatformEventType = "analytics.alert.triggered"
)

type ReportGeneratedEvent struct {
	ReportID    string    `json:"report_id"`
	Name        string    `json:"name"`
	MetricCount int       `json:"metric_count"`
	GeneratedAt time.Time `json:"generated_at"`
}

type DashboardSharedEvent struct {
	DashboardID string    `json:"dashboard_id"`
	Title       string    `json:"title"`
	SharedBy    string    `json:"shared_by"`
	SharedWith  []string  `json:"shared_with"`
	Timestamp   time.Time `json:"timestamp"`
}

type AnomalyDetectedEvent struct {
	MetricName    string  `json:"metric_name"`
	CurrentValue  float64 `json:"current_value"`
	ExpectedValue float64 `json:"expected_value"`
	Deviation     float64 `json:"deviation"`
	DetectedAt    time.Time `json:"detected_at"`
}
