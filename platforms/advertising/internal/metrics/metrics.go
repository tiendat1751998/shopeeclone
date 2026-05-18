package metrics
import ("github.com/prometheus/client_golang/prometheus"; "github.com/prometheus/client_golang/prometheus/promauto")
var (
	CampaignsCreated   = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_ads_campaigns_created_total", Help: "Total campaigns created"})
	AdRequestsTotal    = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_ads_requests_total", Help: "Total ad requests"})
	ImpressionsRecorded = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_ads_impressions_total", Help: "Total impressions"})
	ClicksRecorded     = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_ads_clicks_total", Help: "Total clicks"})
)
