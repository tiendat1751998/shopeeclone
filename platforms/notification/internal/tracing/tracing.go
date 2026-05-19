package tracing

import (
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
)

func Init(serviceName, endpoint string) (func(), error) {
	return observability.InitTracer(serviceName, endpoint)
}
