package creative

import "time"

type CreativeFormat string

const (
	CreativeFormatBanner CreativeFormat = "banner"
	CreativeFormatVideo  CreativeFormat = "video"
	CreativeFormatText   CreativeFormat = "text"
)

type CreativeStatus string

const (
	CreativeStatusDraft        CreativeStatus = "draft"
	CreativeStatusPendingReview CreativeStatus = "pending_review"
	CreativeStatusApproved     CreativeStatus = "approved"
	CreativeStatusRejected     CreativeStatus = "rejected"
)

type CreativeSize struct {
	Width  int
	Height int
	Label  string
}

type CreativePerformance struct {
	Impressions int64
	Clicks      int64
	CTR         float64
}

type Creative struct {
	ID             string
	CampaignID     string
	Name           string
	Format         CreativeFormat
	Status         CreativeStatus
	Content        string
	DestinationURL string
	Sizes          []CreativeSize
	Performance    CreativePerformance
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
