package metrics
import ("github.com/prometheus/client_golang/prometheus"; "github.com/prometheus/client_golang/prometheus/promauto")
var (
	LivestreamsCreated = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_live_created_total", Help: "Total livestreams created"})
	LivestreamsStarted = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_live_started_total", Help: "Total livestreams started"})
	LivestreamsEnded   = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_live_ended_total", Help: "Total livestreams ended"})
	ChatMessagesTotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_live_chat_total", Help: "Total chat messages"})
	ReactionsTotal     = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_live_reactions_total", Help: "Total reactions"})
	GiftsTotal         = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_live_gifts_total", Help: "Total gifts"})
)
