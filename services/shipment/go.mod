module github.com/shopee-clone/shopee/services/shipment

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/prometheus/client_golang v1.19.0
	github.com/redis/go-redis/v9 v9.5.1
	github.com/segmentio/kafka-go v0.4.47
	github.com/shopee-clone/shopee/packages/go-shared v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.6.0
	golang.org/x/crypto v0.24.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
)

replace github.com/shopee-clone/shopee/packages/go-shared => ../../packages/go-shared
