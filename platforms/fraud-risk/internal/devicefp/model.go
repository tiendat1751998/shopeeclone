package devicefp

import "time"

type DeviceProfile struct {
	DeviceID     string    `json:"device_id"`
	UserAgents   []string  `json:"user_agents"`
	IPs          []string  `json:"ips"`
	FirstSeen    time.Time `json:"first_seen"`
	LastSeen     time.Time `json:"last_seen"`
	RiskScore    float64   `json:"risk_score"`
	IsSuspicious bool      `json:"is_suspicious"`
}

type Fingerprint struct {
	DeviceID      string `json:"device_id"`
	UserAgent     string `json:"user_agent"`
	ScreenWidth   int    `json:"screen_width"`
	ScreenHeight  int    `json:"screen_height"`
	ColorDepth    int    `json:"color_depth"`
	Platform      string `json:"platform"`
	Language      string `json:"language"`
	Timezone      string `json:"timezone"`
	FingerprintID string `json:"fingerprint_id"`
	Hash          string `json:"hash"`
}
