package logging

import (
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"go.uber.org/zap"
)

func Init(serviceName, level string) *zap.Logger {
	return observability.InitLogger(serviceName, level)
}
