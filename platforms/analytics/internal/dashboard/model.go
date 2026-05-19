package dashboard

import "time"

type ChartType string

const (
	ChartLine  ChartType = "line"
	ChartBar   ChartType = "bar"
	ChartPie   ChartType = "pie"
	ChartTable ChartType = "table"
	ChartMetric ChartType = "metric"
)

type DataSource struct {
	Type    string                 `json:"type"`
	Query   map[string]interface{} `json:"query,omitempty"`
	Metric  string                 `json:"metric,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type Widget struct {
	ID          string                 `json:"id"`
	DashboardID string                 `json:"dashboard_id"`
	Title       string                 `json:"title"`
	Type        ChartType              `json:"type"`
	Width       int                    `json:"width"`
	Height      int                    `json:"height"`
	PositionX   int                    `json:"position_x"`
	PositionY   int                    `json:"position_y"`
	DataSource  DataSource             `json:"data_source"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Data        interface{}            `json:"data,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type Dashboard struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description,omitempty"`
	OrganizationID string    `json:"organization_id"`
	CreatedBy      string    `json:"created_by"`
	Widgets        []Widget  `json:"widgets,omitempty"`
	IsPublic       bool      `json:"is_public"`
	Tags           []string  `json:"tags,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
