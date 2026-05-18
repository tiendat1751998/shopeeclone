package tracing
import ("github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "github.com/shopee-clone/shopee/platforms/search/internal/config")
func Init(cfg config.OTELConfig) (func(), error) { return observability.InitTracer(cfg.ServiceName, cfg.Endpoint) }
