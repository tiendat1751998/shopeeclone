package cohort

import "time"

type CohortPeriod string

const (
	CohortDay   CohortPeriod = "day"
	CohortWeek  CohortPeriod = "week"
	CohortMonth CohortPeriod = "month"
)

type CohortDefinition struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Period           CohortPeriod `json:"period"`
	AcquisitionField string       `json:"acquisition_field"`
	TimeRange        string       `json:"time_range"`
	CreatedAt        time.Time    `json:"created_at"`
}

type CohortAnalysis struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Period      CohortPeriod        `json:"period"`
	PeriodLabel string              `json:"period_label"`
	Cohorts     []CohortRow         `json:"cohorts"`
	Matrix      [][]CohortCell      `json:"matrix"`
	Periods     []int               `json:"periods"`
	AnalyzedAt  time.Time           `json:"analyzed_at"`
}

type CohortRow struct {
	ID               string  `json:"id"`
	PeriodStart      string  `json:"period_start"`
	UserCount        int64   `json:"user_count"`
	Retention        []float64 `json:"retention"`
	AcquisitionSource string `json:"acquisition_source,omitempty"`
}

type CohortCell struct {
	UserCount     int64   `json:"user_count"`
	RetentionRate float64 `json:"retention_rate"`
}

type RetentionPoint struct {
	PeriodOffset int     `json:"period_offset"`
	UserCount    int64   `json:"user_count"`
	Rate         float64 `json:"rate"`
}
