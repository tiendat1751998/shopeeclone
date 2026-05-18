package metrics
import ("github.com/prometheus/client_golang/prometheus"; "github.com/prometheus/client_golang/prometheus/promauto")
var (
	FraudScoresComputed = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_fraud_scores_computed_total", Help: "Total fraud scores computed"})
	FraudCasesCreated   = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_fraud_cases_created_total", Help: "Total fraud cases created"})
	FraudBlocks         = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_fraud_blocks_total", Help: "Total fraud blocks"})
)
