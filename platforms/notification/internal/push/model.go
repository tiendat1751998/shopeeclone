package push

import "time"

type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformWeb     Platform = "web"
)

type PushDevice struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	Platform  Platform  `json:"platform"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PushToken struct {
	Token    string   `json:"token"`
	Platform Platform `json:"platform"`
}

type PushNotificationRequest struct {
	UserID string            `json:"user_id"`
	Title  string            `json:"title"`
	Body   string            `json:"body"`
	Data   map[string]string `json:"data,omitempty"`
	Badge  int               `json:"badge,omitempty"`
	Sound  string            `json:"sound,omitempty"`
}

type BulkPushRequest struct {
	UserIDs []string          `json:"user_ids"`
	Title   string            `json:"title"`
	Body    string            `json:"body"`
	Data    map[string]string `json:"data,omitempty"`
}

type PushResult struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Error   string `json:"error,omitempty"`
}
