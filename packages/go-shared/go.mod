module github.com/shopee-clone/shopee/packages/go-shared

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis_rate/v10 v10.0.1
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/redis/go-redis/v9 v9.5.1
	github.com/segmentio/kafka-go v0.4.47
	github.com/sony/gobreaker v0.5.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.50.0
	go.opentelemetry.io/otel v1.25.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.25.0
	go.opentelemetry.io/otel/sdk v1.25.0
	go.opentelemetry.io/otel/trace v1.25.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.6.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
)
