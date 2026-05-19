package analytics

import "time"

type Impression struct {
	ID         string
	CampaignID string
	CreativeID string
	UserID     string
	Timestamp  time.Time
	Cost       float64
	Device     string
	Location   string
}

type Click struct {
	ID           string
	ImpressionID string
	CampaignID   string
	CreativeID   string
	UserID       string
	Timestamp    time.Time
	Cost         float64
}

type Conversion struct {
	ID             string
	ClickID        string
	CampaignID     string
	CreativeID     string
	UserID         string
	Timestamp      time.Time
	Revenue        float64
	ConversionType string
}

type AnalyticsReport struct {
	CampaignID  string
	Impressions int64
	Clicks      int64
	Conversions int64
	Spend       float64
	Revenue     float64
	CTR         float64
	CVR         float64
	CPC         float64
	CPM         float64
	ROAS        float64
}

type ReportFilter struct {
	CampaignID   string
	CreativeID   string
	StartDate    string
	EndDate      string
}
