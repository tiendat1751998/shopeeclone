package analytics

import "time"

type AggregationType string

const (
	AggSum           AggregationType = "sum"
	AggCount         AggregationType = "count"
	AggAvg           AggregationType = "avg"
	AggMin           AggregationType = "min"
	AggMax           AggregationType = "max"
	AggDistinctCount AggregationType = "distinct_count"
)

type MetricType string

const (
	MetricTotalUsers       MetricType = "total_users"
	MetricActiveUsers      MetricType = "active_users"
	MetricRevenue          MetricType = "revenue"
	MetricOrders           MetricType = "orders"
	MetricAOV              MetricType = "aov"
	MetricConversionRate   MetricType = "conversion_rate"
	MetricPageviews        MetricType = "pageviews"
	MetricSessions         MetricType = "sessions"
	MetricBounceRate       MetricType = "bounce_rate"
	MetricAvgSessionDuration MetricType = "avg_session_duration"
)

type DimensionType string

const (
	DimDate     DimensionType = "date"
	DimHour     DimensionType = "hour"
	DimSource   DimensionType = "source"
	DimDevice   DimensionType = "device"
	DimCountry  DimensionType = "country"
	DimCampaign DimensionType = "campaign"
	DimEventType DimensionType = "event_type"
)

type TimeRangeType string

const (
	TimeToday     TimeRangeType = "today"
	TimeYesterday TimeRangeType = "yesterday"
	TimeLast7d    TimeRangeType = "last_7d"
	TimeLast30d   TimeRangeType = "last_30d"
	TimeLast90d   TimeRangeType = "last_90d"
	TimeCustom    TimeRangeType = "custom"
)

type TimeRange struct {
	Type    TimeRangeType `json:"type"`
	StartAt *time.Time    `json:"start_at,omitempty"`
	EndAt   *time.Time    `json:"end_at,omitempty"`
}

type Metric struct {
	Name       MetricType      `json:"name"`
	Aggregation AggregationType `json:"aggregation"`
	Alias      string           `json:"alias,omitempty"`
}

type QueryFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type AnalyticsQuery struct {
	Metrics     []Metric          `json:"metrics"`
	Dimensions  []DimensionType   `json:"dimensions,omitempty"`
	TimeRange   TimeRange         `json:"time_range"`
	Filters     []QueryFilter     `json:"filters,omitempty"`
	GroupBy     []string          `json:"group_by,omitempty"`
	OrderBy     string            `json:"order_by,omitempty"`
	OrderDir    string            `json:"order_dir,omitempty"`
	Limit       int               `json:"limit,omitempty"`
	Offset      int               `json:"offset,omitempty"`
}

type QueryResult struct {
	Columns []string         `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
	Total   int64             `json:"total"`
}

type MetricValue struct {
	Name  MetricType      `json:"name"`
	Value float64         `json:"value"`
	Label string          `json:"label,omitempty"`
}

type Report struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Query          AnalyticsQuery `json:"query"`
	Result         *QueryResult   `json:"result,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	OrganizationID string         `json:"organization_id"`
}
