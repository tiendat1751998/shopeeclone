package preferences

import "time"

type ChannelOptIn struct {
	Push  bool `json:"push"`
	Email bool `json:"email"`
	SMS   bool `json:"sms"`
	InApp bool `json:"inapp"`
}

type CategoryPreferences map[string]bool

type QuietHours struct {
	Enabled bool   `json:"enabled"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Timezone string `json:"timezone"`
}

type UserPreference struct {
	UserID          string               `json:"user_id"`
	ChannelOptIn    ChannelOptIn         `json:"channel_opt_in"`
	Categories      CategoryPreferences  `json:"categories"`
	QuietHours      QuietHours           `json:"quiet_hours"`
	EmailDigest     bool                 `json:"email_digest"`
	PushEnabled     bool                 `json:"push_enabled"`
	SMSPromotions   bool                 `json:"sms_promotions"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

type SuppressionEntry struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdatePreferenceRequest struct {
	ChannelOptIn  *ChannelOptIn        `json:"channel_opt_in,omitempty"`
	Categories    *CategoryPreferences `json:"categories,omitempty"`
	QuietHours    *QuietHours          `json:"quiet_hours,omitempty"`
	EmailDigest   *bool                `json:"email_digest,omitempty"`
	PushEnabled   *bool                `json:"push_enabled,omitempty"`
	SMSPromotions *bool                `json:"sms_promotions,omitempty"`
}
