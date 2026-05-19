package transactionmon

import "time"

type TransactionMonitor struct {
	UserID            string    `json:"user_id"`
	DailyCount        int       `json:"daily_count"`
	HourlyCount       int       `json:"hourly_count"`
	DailyVolume       float64   `json:"daily_volume"`
	HourlyVolume      float64   `json:"hourly_volume"`
	AvgTicket         float64   `json:"avg_ticket"`
	LastTransactions  []float64 `json:"last_transactions"`
	LastDailyReset    time.Time `json:"last_daily_reset"`
	LastHourlyReset   time.Time `json:"last_hourly_reset"`
}

type TransactionRecord struct {
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Location  string    `json:"location"`
	IP        string    `json:"ip"`
	DeviceID  string    `json:"device_id"`
}

type AnomalyResult struct {
	HasAnomaly      bool     `json:"has_anomaly"`
	Reasons         []string `json:"reasons"`
	CurrentVelocity int      `json:"current_velocity"`
	AvgTicket       float64  `json:"avg_ticket"`
}
