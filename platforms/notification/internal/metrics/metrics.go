package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	NotificationsSentTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "notifications_sent_total",
		Help: "Total notifications sent by channel",
	}, []string{"channel"})

	NotificationsFailedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "notifications_failed_total",
		Help: "Total notifications failed by channel",
	}, []string{"channel"})

	PushSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "push_sent_total",
		Help: "Total push notifications sent",
	})

	PushFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "push_failed_total",
		Help: "Total push notifications failed",
	})

	EmailSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "email_sent_total",
		Help: "Total emails sent",
	})

	EmailFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "email_failed_total",
		Help: "Total emails failed",
	})

	EmailBouncedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "email_bounced_total",
		Help: "Total emails bounced",
	})

	SMSSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sms_sent_total",
		Help: "Total SMS sent",
	})

	SMSFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sms_failed_total",
		Help: "Total SMS failed",
	})

	InAppSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "inapp_sent_total",
		Help: "Total in-app notifications sent",
	})

	NotificationLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "notification_latency_seconds",
		Help:    "Notification delivery latency",
		Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	}, []string{"channel"})

	RateLimitHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notification_rate_limit_hits_total", Help: "Total rate limit hits",
	})

	NotificationsSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notifications_sent_count", Help: "Total notifications sent (bare count)",
	})
)
