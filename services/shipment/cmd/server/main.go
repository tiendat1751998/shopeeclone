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
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	sharedRedis "github.com/shopee-clone/shopee/packages/go-shared/pkg/redis"
	"github.com/shopee-clone/shopee/services/shipment/internal/application"
	"github.com/shopee-clone/shopee/services/shipment/internal/config"
	"github.com/shopee-clone/shopee/services/shipment/internal/health"
	"github.com/shopee-clone/shopee/services/shipment/internal/infrastructure/kafka"
	"github.com/shopee-clone/shopee/services/shipment/internal/infrastructure/mysql"
	redisinfra "github.com/shopee-clone/shopee/services/shipment/internal/infrastructure/redis"
	"github.com/shopee-clone/shopee/services/shipment/internal/tracing"
	grpctransport "github.com/shopee-clone/shopee/services/shipment/internal/transport/grpc"
	httptransport "github.com/shopee-clone/shopee/services/shipment/internal/transport/http"
	"github.com/shopee-clone/shopee/services/shipment/internal/transport/http/middleware"
	pb "github.com/shopee-clone/shopee/services/shipment/proto/shipment/v1"
	"go.uber.org/zap"
	automaxprocs "go.uber.org/automaxprocs/maxprocs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var version = "1.0.0"

func init() {
	// Tune GC for low-latency: more frequent GCs, less heap growth
	if gogc := os.Getenv("GOGC"); gogc == "" {
		os.Setenv("GOGC", "50")
	}
}

func main() {
	cfg := config.Load()

	logger := observability.InitLogger(cfg.AppName, cfg.LogLevel)

	// Auto-tune GOMAXPROCS for container environments
	if _, err := automaxprocs.Set(); err != nil {
		logger.Warn("failed to set automaxprocs", zap.Error(err))
	}
	shutdownTracer, err := tracing.Init(cfg.OpenTelemetry)
	if err != nil { logger.Fatal("failed to init tracer", zap.Error(err)) }
	defer shutdownTracer(); defer observability.Sync()

	db, err := mysql.NewDB(cfg.MySQL)
	if err != nil { logger.Fatal("failed to connect to mysql", zap.Error(err)) }
	defer db.Close()

	redisClient, err := sharedRedis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil { logger.Warn("redis not available", zap.Error(err)); redisClient = nil }

	shipmentRepo := mysql.NewShipmentRepository(db)
	redisStore := redisinfra.NewStore(redisClient, cfg.Redis)
	var kafkaProducer *kafka.Producer
	if len(cfg.Kafka.Brokers) > 0 && cfg.Kafka.Brokers[0] != "" { kafkaProducer = kafka.NewProducer(cfg.Kafka) }
	shipmentService := application.NewShipmentService(cfg, shipmentRepo, redisStore, kafkaProducer)

	gin.SetMode(getGinMode(cfg.AppEnv))
	engine := gin.New()
	var authMw gin.HandlerFunc
	if cfg.JWT.AccessSecret != "" && cfg.JWT.AccessSecret != "change-me-in-production" { authMw = middleware.JWTAuth(cfg.JWT) }
	handler := httptransport.NewHandler(shipmentService)
	router := httptransport.NewRouter(handler, authMw)
	router.Setup(engine)

	healthChecker := health.NewChecker(cfg.AppName, version, db, redisClient)
	engine.GET("/health/live", healthChecker.LivenessHandler())
	engine.GET("/health/ready", healthChecker.ReadinessHandler())
	engine.GET("/health", healthChecker.ReadinessHandler())

	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", cfg.HTTPPort), Handler: engine, ReadTimeout:       5 * time.Second, WriteTimeout:      10 * time.Second, IdleTimeout:       120 * time.Second}
	grpcServer := grpc.NewServer()
	pb.RegisterShipmentServiceServer(grpcServer, grpctransport.NewShipmentGRPCServer(shipmentService))
	grpc_health_v1.RegisterHealthServer(grpcServer, &grpcHealthServer{})
	reflection.Register(grpcServer)
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil { logger.Fatal("failed to listen grpc", zap.Error(err)) }

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-quit
		cancel()
	}()

	go func() {
		logger.Info("starting shipment service", zap.Int("http_port", cfg.HTTPPort), zap.Int("grpc_port", cfg.GRPCPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed { logger.Fatal("http server failed", zap.Error(err)) }
	}()
	go func() { if err := grpcServer.Serve(grpcListener); err != nil { logger.Fatal("grpc server failed", zap.Error(err)) } }()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("panic in shipment outbox worker", zap.Any("recover", r))
			}
		}()
		ticker := time.NewTicker(5 * time.Second); defer ticker.Stop()
		for { select { case <-ticker.C: func(){ workerCtx, workerCancel := context.WithTimeout(context.Background(), 30*time.Second); defer workerCancel(); shipmentService.ProcessOutboxEvents(workerCtx) }(); case <-quit: return } }
	}()

	<-quit; logger.Info("shutting down shipment service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second); defer shutdownCancel()
	grpcServer.GracefulStop(); httpServer.Shutdown(shutdownCtx)
	if redisClient != nil { redisClient.Close() }
	if kafkaProducer != nil { kafkaProducer.Close() }
	logger.Info("shipment service stopped")
}

type grpcHealthServer struct{}
func (s *grpcHealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}
func (s *grpcHealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return stream.Send(&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING})
}
func getGinMode(env string) string { if env == "production" || env == "staging" { return gin.ReleaseMode }; return gin.DebugMode }
