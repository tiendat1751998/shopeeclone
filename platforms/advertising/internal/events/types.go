package events

import "time"

type EventType string

const (
	EventCampaignCreated     EventType = "campaign.created"
	EventAuctionWon          EventType = "auction.won"
	EventImpressionRecorded  EventType = "impression.recorded"
	EventClickRecorded       EventType = "click.recorded"
	EventConversionRecorded  EventType = "conversion.recorded"
)

type CampaignCreated struct {
	CampaignID string    `json:"campaign_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Budget     float64   `json:"budget"`
	Timestamp  time.Time `json:"timestamp"`
}

type AuctionWon struct {
	CampaignID  string    `json:"campaign_id"`
	CreativeID  string    `json:"creative_id"`
	UserID      string    `json:"user_id"`
	BidAmount   float64   `json:"bid_amount"`
	SecondPrice float64   `json:"second_price"`
	Timestamp   time.Time `json:"timestamp"`
}

type ImpressionRecorded struct {
	ImpressionID string    `json:"impression_id"`
	CampaignID   string    `json:"campaign_id"`
	CreativeID   string    `json:"creative_id"`
	UserID       string    `json:"user_id"`
	Cost         float64   `json:"cost"`
	Timestamp    time.Time `json:"timestamp"`
}

type ClickRecorded struct {
	ClickID      string    `json:"click_id"`
	ImpressionID string    `json:"impression_id"`
	CampaignID   string    `json:"campaign_id"`
	CreativeID   string    `json:"creative_id"`
	UserID       string    `json:"user_id"`
	Cost         float64   `json:"cost"`
	Timestamp    time.Time `json:"timestamp"`
}
