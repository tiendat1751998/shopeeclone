package campaign

import "time"

type CampaignStatus string

const (
	CampaignStatusDraft  CampaignStatus = "draft"
	CampaignStatusActive CampaignStatus = "active"
	CampaignStatusPaused CampaignStatus = "paused"
	CampaignStatusEnded  CampaignStatus = "ended"
)

type CampaignType string

const (
	CampaignTypeCPC CampaignType = "CPC"
	CampaignTypeCPM CampaignType = "CPM"
	CampaignTypeCPA CampaignType = "CPA"
)

type Budget struct {
	Daily    float64
	Lifetime float64
}

type DateRange struct {
	Start time.Time
	End   time.Time
}

type Targeting struct {
	Demographics *Demographic
	Interests    []string
	Locations    []string
	Devices      []string
}

type Demographic struct {
	MinAge int
	MaxAge int
	Gender string
}

type Campaign struct {
	ID           string
	Name         string
	Status       CampaignStatus
	Type         CampaignType
	Budget       Budget
	DateRange    DateRange
	Targeting    Targeting
	CreativeIDs  []string
	BidAmount    float64
	TargetCPA    float64
	QualityScore float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
