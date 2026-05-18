package logging
import ("github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.uber.org/zap")
func Init(s, l string) *zap.Logger { return observability.InitLogger(s, l) }
