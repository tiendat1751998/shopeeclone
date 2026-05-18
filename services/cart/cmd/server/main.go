package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sharedRedis "github.com/shopee-clone/shopee/packages/go-shared/pkg/redis"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/services/cart/internal/application"
	"github.com/shopee-clone/shopee/services/cart/internal/config"
	"github.com/shopee-clone/shopee/services/cart/internal/infrastructure/mysql"
	redisinfra "github.com/shopee-clone/shopee/services/cart/internal/infrastructure/redis"
	httptransport "github.com/shopee-clone/shopee/services/cart/internal/transport/http"
	kafkatransport "github.com/shopee-clone/shopee/services/cart/internal/transport/kafka"
	"go.uber.org/zap"
)

var version = "1.0.0"

func main() {
	cfg := config.Load()
	logger := observability.InitLogger(cfg.AppName, cfg.LogLevel)

	shutdownTracer, err := observability.InitTracer(cfg.OpenTelemetry.ServiceName, cfg.OpenTelemetry.Endpoint)
	if err != nil {
		logger.Fatal("failed to init tracer", zap.Error(err))
	}
	defer shutdownTracer()
	defer observability.Sync()

	db, err := mysql.NewDB(cfg.MySQL)
	if err != nil {
		logger.Fatal("failed to connect to mysql", zap.Error(err))
	}
	defer db.Close()

	redisClient, err := sharedRedis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Warn("redis not available", zap.Error(err))
		redisClient = nil
	}
	redisStore := redisinfra.NewStore(redisClient, cfg.Redis)

	cartRepo := mysql.NewCartRepository(db)
	itemRepo := mysql.NewCartItemRepository(db)
	snapshotRepo := mysql.NewCartSnapshotRepository(db)
	mergeRepo := mysql.NewCartMergeHistoryRepository(db)

	var publisher application.EventPublisher
	if len(cfg.Kafka.Brokers) > 0 && cfg.Kafka.Brokers[0] != "" {
		producer := kafkatransport.NewProducer(cfg.Kafka.Brokers, cfg.AppName)
		defer producer.Close()
		publisher = producer
	}

	cartService := application.NewCartService(
		cartRepo, itemRepo, snapshotRepo, mergeRepo,
		redisStore, cfg.CartTTL, cfg.MaxCartItems, cfg.MaxQuantityPerItem, publisher,
	)

	healthChecker := health.NewChecker(cfg.AppName, version)
	healthChecker.AddCheck("database", func(ctx context.Context) error { return db.Ping() })
	if redisClient != nil {
		healthChecker.AddCheck("redis", func(ctx context.Context) error { return redisClient.Ping(ctx).Err() })
	}

	gin.SetMode(getGinMode(cfg.AppEnv))
	engine := gin.New()

	handler := httptransport.NewHandler(cartService)
	httpRouter := httptransport.NewRouter(handler, healthChecker)
	httpRouter.Setup(engine)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("starting cart service", zap.Int("http_port", cfg.HTTPPort), zap.String("env", cfg.AppEnv))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("shutting down cart service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown error", zap.Error(err))
	}
	if redisClient != nil {
		redisClient.Close()
	}

	logger.Info("cart service stopped")
}

func getGinMode(env string) string {
	switch env {
	case "production", "staging":
		return gin.ReleaseMode
	default:
		return gin.DebugMode
	}
}
