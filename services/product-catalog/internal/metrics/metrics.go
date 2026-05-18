package metrics
import ("github.com/prometheus/client_golang/prometheus"; "github.com/prometheus/client_golang/prometheus/promauto")
var (
	ProductsCreated = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_catalog_products_created_total", Help: "Total products created"})
	ProductsUpdated = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_catalog_products_updated_total", Help: "Total products updated"})
	SKUsCreated     = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_catalog_skus_created_total", Help: "Total SKUs created"})
	IdempotentRequests = promauto.NewCounter(prometheus.CounterOpts{Name: "shopee_catalog_idempotent_requests_total", Help: "Idempotent requests"})
)
