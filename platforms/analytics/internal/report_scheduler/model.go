package report_scheduler

import "time"

type ScheduleFrequency string

const (
	FreqDaily   ScheduleFrequency = "daily"
	FreqWeekly  ScheduleFrequency = "weekly"
	FreqMonthly ScheduleFrequency = "monthly"
)

type DeliveryChannel string

const (
	ChannelEmail    DeliveryChannel = "email"
	ChannelDownload DeliveryChannel = "download"
	ChannelWebhook  DeliveryChannel = "webhook"
)

type ScheduledReport struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	Query            map[string]interface{} `json:"query"`
	Frequency        ScheduleFrequency `json:"frequency"`
	DeliveryChannel  DeliveryChannel   `json:"delivery_channel"`
	Recipients       []string          `json:"recipients,omitempty"`
	WebhookURL       string            `json:"webhook_url,omitempty"`
	Format           string            `json:"format"`
	TimeZone         string            `json:"time_zone"`
	NextRunAt        time.Time         `json:"next_run_at"`
	LastRunAt        *time.Time        `json:"last_run_at,omitempty"`
	LastStatus       string            `json:"last_status,omitempty"`
	IsActive         bool              `json:"is_active"`
	CreatedBy        string            `json:"created_by"`
	OrganizationID   string            `json:"organization_id"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

type ReportGeneration struct {
	ID         string    `json:"id"`
	ReportID   string    `json:"report_id"`
	Status     string    `json:"status"`
	Data       interface{} `json:"data,omitempty"`
	Error      string    `json:"error,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
}
