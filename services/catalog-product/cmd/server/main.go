package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sharedConfig "github.com/shopee-clone/shopee/packages/go-shared/pkg/config"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/kafka"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	sharedRedis "github.com/shopee-clone/shopee/packages/go-shared/pkg/redis"
	catalogGrpc "github.com/shopee-clone/shopee/services/catalog-product/internal/delivery/grpc"
	catalogHttp "github.com/shopee-clone/shopee/services/catalog-product/internal/delivery/http"
	"github.com/shopee-clone/shopee/services/catalog-product/internal/repository"
	"github.com/shopee-clone/shopee/services/catalog-product/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var version = "0.1.0"

func main() {
	cfg := sharedConfig.Load()
	logger := observability.InitLogger("catalog-product", cfg.LogLevel)

	shutdownTracer, err := observability.InitTracer("catalog-product", cfg.OpenTelemetry.Endpoint)
	if err != nil {
		logger.Fatal("failed to init tracer", zap.Error(err))
	}
	defer shutdownTracer()
	defer observability.Sync()

	mongoClient, err := repository.NewMongoClient(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		logger.Fatal("failed to connect to mongodb", zap.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		mongoClient.Disconnect(ctx)
	}()

	redisClient, err := sharedRedis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Warn("redis not available, continuing without cache", zap.Error(err))
		redisClient = nil
	}

	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers, "catalog-product")
	defer kafkaProducer.Close()

	productRepo := repository.NewProductRepository(mongoClient, cfg.MongoDB.Database)
	categoryRepo := repository.NewCategoryRepository(mongoClient, cfg.MongoDB.Database)
	productCache := repository.NewProductCache(redisClient)
	productUseCase := usecase.NewProductUseCase(productRepo, productCache, kafkaProducer)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo, kafkaProducer)

	healthChecker := health.NewChecker("catalog-product", version)
	healthChecker.AddCheck("mongodb", func(ctx context.Context) error {
		return mongoClient.Ping(ctx, nil)
	})
	if redisClient != nil {
		healthChecker.AddCheck("redis", func(ctx context.Context) error {
			return redisClient.Ping(ctx).Err()
		})
	}

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())
	router.Use(middleware.OTelMiddleware("catalog-product"))
	router.Use(observability.ObserveHTTPMetrics("catalog-product"))

	router.GET("/health", healthChecker.LivenessHandler())
	router.GET("/ready", healthChecker.ReadinessHandler())
	router.GET("/metrics", observability.MetricsHandler())

	httpHandler := catalogHttp.NewHandler(productUseCase, categoryUseCase)
	api := router.Group("/api/v1")
	{
		products := api.Group("/products")
		{
			products.GET("", httpHandler.ListProducts)
			products.GET("/:spu_id", httpHandler.GetProduct)
			products.POST("", httpHandler.CreateProduct)
			products.PUT("/:spu_id", httpHandler.UpdateProduct)
			products.DELETE("/:spu_id", httpHandler.DeleteProduct)
		}
		categories := api.Group("/categories")
		{
			categories.GET("", httpHandler.ListCategories)
			categories.GET("/:category_id", httpHandler.GetCategory)
		}
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(catalogGrpc.OTelUnaryServerInterceptor()),
	)
	grpc_health_v1.RegisterHealthServer(grpcServer, &catalogGrpc.HealthServer{})
	reflection.Register(grpcServer)

	grpcPort := cfg.Port + 1
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Fatal("failed to listen grpc", zap.Int("port", grpcPort), zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("starting grpc server", zap.Int("port", grpcPort))
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("grpc server failed", zap.Error(err))
		}
	}()

	go func() {
		logger.Info("starting http server", zap.Int("port", cfg.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown error", zap.Error(err))
	}

	logger.Info("server stopped")
}
