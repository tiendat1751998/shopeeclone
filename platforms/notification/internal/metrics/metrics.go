package metrics
import ("github.com/prometheus/client_golang/prometheus"; "github.com/prometheus/client_golang/prometheus/promauto")
var (
	NotificationsSent = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_notification_sent_total", Help: "Total notifications sent"})
	RateLimitHits     = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_notification_rate_limit_hits_total", Help: "Rate limit hits"})
	DeliveryErrors    = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_notification_delivery_errors_total", Help: "Delivery errors"})
	NotificationLatency = promauto.NewHistogram(prometheus.HistogramOpts{Name: "shopee_notification_latency_seconds", Help: "Notification latency", Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1}})
)
