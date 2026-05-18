module github.com/shopee-clone/shopee/services/catalog-product

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.6.0
	github.com/redis/go-redis/v9 v9.5.1
	github.com/shopee-clone/shopee/packages/go-shared v0.0.0
	go.mongodb.org/mongo-driver v1.15.0
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.50.0
	go.opentelemetry.io/otel v1.25.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
)

replace github.com/shopee-clone/shopee/packages/go-shared => ../../packages/go-shared
