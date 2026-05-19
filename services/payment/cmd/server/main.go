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
	"github.com/shopee-clone/shopee/services/payment/internal/application"
	"github.com/shopee-clone/shopee/services/payment/internal/config"
	"github.com/shopee-clone/shopee/services/payment/internal/health"
	"github.com/shopee-clone/shopee/services/payment/internal/infrastructure/kafka"
	"github.com/shopee-clone/shopee/services/payment/internal/infrastructure/mysql"
	redisinfra "github.com/shopee-clone/shopee/services/payment/internal/infrastructure/redis"
	"github.com/shopee-clone/shopee/services/payment/internal/tracing"
	grpctransport "github.com/shopee-clone/shopee/services/payment/internal/transport/grpc"
	httptransport "github.com/shopee-clone/shopee/services/payment/internal/transport/http"
	"github.com/shopee-clone/shopee/services/payment/internal/transport/http/middleware"
	pb "github.com/shopee-clone/shopee/services/payment/proto/payment/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var version = "1.0.0"

func main() {
	cfg := config.Load()
	logger := observability.InitLogger(cfg.AppName, cfg.LogLevel)

	shutdownTracer, err := tracing.Init(cfg.OpenTelemetry)
	if err != nil { logger.Fatal("failed to init tracer", zap.Error(err)) }
	defer shutdownTracer()
	defer observability.Sync()

	db, err := mysql.NewDB(cfg.MySQL)
	if err != nil { logger.Fatal("failed to connect to mysql", zap.Error(err)) }
	defer db.Close()

	redisClient, err := sharedRedis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil { logger.Warn("redis not available", zap.Error(err)); redisClient = nil }

	paymentRepo := mysql.NewPaymentRepository(db)
	redisStore := redisinfra.NewStore(redisClient, cfg.Redis)

	var kafkaProducer *kafka.Producer
	if len(cfg.Kafka.Brokers) > 0 && cfg.Kafka.Brokers[0] != "" {
		kafkaProducer = kafka.NewProducer(cfg.Kafka)
	}

	paymentService := application.NewPaymentService(cfg, paymentRepo, redisStore, kafkaProducer)

	gin.SetMode(getGinMode(cfg.AppEnv))
	engine := gin.New()

	var authMw gin.HandlerFunc
	if cfg.JWT.AccessSecret != "" && cfg.JWT.AccessSecret != "change-me-in-production" {
		authMw = middleware.JWTAuth(cfg.JWT)
	}

	handler := httptransport.NewHandler(paymentService)
	router := httptransport.NewRouter(handler, authMw)
	router.Setup(engine)

	healthChecker := health.NewChecker(cfg.AppName, version, db, redisClient)
	engine.GET("/health/live", healthChecker.LivenessHandler())
	engine.GET("/health/ready", healthChecker.ReadinessHandler())

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.HTTPPort), Handler: engine,
		ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second,
	}

	grpcServer := grpc.NewServer()
	grpcPaymentServer := grpctransport.NewPaymentGRPCServer(paymentService)
	pb.RegisterPaymentServiceServer(grpcServer, grpcPaymentServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, &grpcHealthServer{})
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil { logger.Fatal("failed to listen grpc", zap.Error(err)) }

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("starting payment service", zap.Int("http_port", cfg.HTTPPort), zap.Int("grpc_port", cfg.GRPCPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("grpc server failed", zap.Error(err))
		}
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				func() {
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					if err := paymentService.ProcessOutboxEvents(ctx); err != nil {
						zap.L().Warn("outbox processing failed", zap.Error(err))
					}
				}()
			case <-quit:
				return
			}
		}
	}()

	<-quit
	logger.Info("shutting down payment service...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	httpServer.Shutdown(shutdownCtx)
	if redisClient != nil { redisClient.Close() }
	if kafkaProducer != nil { kafkaProducer.Close() }
	logger.Info("payment service stopped")
}

type grpcHealthServer struct{}

func (s *grpcHealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (s *grpcHealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return stream.Send(&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING})
}

func getGinMode(env string) string {
	switch env {
	case "production", "staging": return gin.ReleaseMode
	default: return gin.DebugMode
	}
}
