package campaign

import "time"

type CampaignType string

const (
	TypePromotional    CampaignType = "promotional"
	TypeTransactional  CampaignType = "transactional"
	TypeOnboarding     CampaignType = "onboarding"
	TypeReEngagement   CampaignType = "re_engagement"
)

type Channel string

const (
	ChannelPush Channel = "push"
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelInApp Channel = "inapp"
)

type Status string

const (
	StatusDraft      Status = "draft"
	StatusScheduled  Status = "scheduled"
	StatusRunning    Status = "running"
	StatusPaused     Status = "paused"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
)

var ValidTransitions = map[Status][]Status{
	StatusDraft:     {StatusScheduled, StatusRunning, StatusCancelled},
	StatusScheduled: {StatusRunning, StatusCancelled},
	StatusRunning:   {StatusPaused, StatusCompleted},
	StatusPaused:    {StatusRunning, StatusCancelled},
}

type Schedule struct {
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
	Timezone string    `json:"timezone"`
}

type Campaign struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Type            CampaignType `json:"type"`
	Channel         Channel      `json:"channel"`
	Status          Status       `json:"status"`
	Schedule        Schedule     `json:"schedule"`
	AudienceQuery   string       `json:"audience_query"`
	ContentTemplate string       `json:"content_template"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type CreateCampaignRequest struct {
	Name            string       `json:"name"`
	Type            CampaignType `json:"type"`
	Channel         Channel      `json:"channel"`
	Schedule        Schedule     `json:"schedule"`
	AudienceQuery   string       `json:"audience_query"`
	ContentTemplate string       `json:"content_template"`
}

type UpdateCampaignRequest struct {
	Name            *string       `json:"name,omitempty"`
	Type            *CampaignType `json:"type,omitempty"`
	Channel         *Channel      `json:"channel,omitempty"`
	Schedule        *Schedule     `json:"schedule,omitempty"`
	AudienceQuery   *string       `json:"audience_query,omitempty"`
	ContentTemplate *string       `json:"content_template,omitempty"`
}
