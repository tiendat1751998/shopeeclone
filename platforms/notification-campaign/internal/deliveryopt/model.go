package deliveryopt

import "time"

type UserEngagementPattern struct {
	UserID        string    `json:"user_id"`
	Channel       string    `json:"channel"`
	PeakOpenHour  int       `json:"peak_open_hour"`
	PeakClickHour int       `json:"peak_click_hour"`
	OpenRate      float64   `json:"open_rate"`
	ClickRate     float64   `json:"click_rate"`
	SampleSize    int       `json:"sample_size"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SendTimeOptimization struct {
	UserID     string `json:"user_id"`
	Channel    string `json:"channel"`
	BestHour   int    `json:"best_hour"`
	Confidence string `json:"confidence"`
}

type ThrottleConfig struct {
	CampaignMessagesPerHour int               `json:"campaign_messages_per_hour"`
	ChannelMessagesPerHour  map[string]int    `json:"channel_messages_per_hour"`
}

type ChannelFallbackPlan struct {
	Channels []string `json:"channels"`
}

var DefaultFallbackOrder = []string{"push", "email", "sms", "inapp"}

type PriorityLevel int

const (
	PriorityBulk          PriorityLevel = 0
	PriorityPromotional   PriorityLevel = 1
	PriorityTransactional PriorityLevel = 2
	PriorityCritical      PriorityLevel = 3
)

type QueuedMessage struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Channel     string        `json:"channel"`
	Content     string        `json:"content"`
	Priority    PriorityLevel `json:"priority"`
	CampaignID  string        `json:"campaign_id"`
	CreatedAt   time.Time     `json:"created_at"`
}

type SendRequest struct {
	UserID    string                 `json:"user_id"`
	Channel   string                 `json:"channel"`
	Subject   string                 `json:"subject"`
	Body      string                 `json:"body"`
	Priority  PriorityLevel          `json:"priority"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type SendResult struct {
	Success     bool   `json:"success"`
	FinalChannel string `json:"final_channel"`
	Error       string `json:"error,omitempty"`
}
