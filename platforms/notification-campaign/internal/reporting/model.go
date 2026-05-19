package reporting

import "time"

type CampaignReport struct {
	CampaignID      string    `json:"campaign_id"`
	SentCount       int       `json:"sent_count"`
	DeliveredCount  int       `json:"delivered_count"`
	OpenedCount     int       `json:"opened_count"`
	ClickedCount    int       `json:"clicked_count"`
	ConvertedCount  int       `json:"converted_count"`
	BouncedCount    int       `json:"bounced_count"`
	UnsubscribedCount int     `json:"unsubscribed_count"`
	RevenueAttributed float64 `json:"revenue_attributed"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type AggregatedReport struct {
	TotalCampaigns     int     `json:"total_campaigns"`
	TotalSent          int     `json:"total_sent"`
	TotalDelivered     int     `json:"total_delivered"`
	TotalOpened        int     `json:"total_opened"`
	TotalClicked       int     `json:"total_clicked"`
	TotalConverted     int     `json:"total_converted"`
	TotalBounced       int     `json:"total_bounced"`
	TotalUnsubscribed  int     `json:"total_unsubscribed"`
	TotalRevenue       float64 `json:"total_revenue"`
	OverallOpenRate    float64 `json:"overall_open_rate"`
	OverallClickRate   float64 `json:"overall_click_rate"`
	OverallConversionRate float64 `json:"overall_conversion_rate"`
}

type TrackEventRequest struct {
	CampaignID string  `json:"campaign_id"`
	UserID     string  `json:"user_id"`
	Revenue    float64 `json:"revenue,omitempty"`
}
