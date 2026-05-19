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
	"github.com/shopee-clone/shopee/services/auth/internal/application"
	"github.com/shopee-clone/shopee/services/auth/internal/config"
	"github.com/shopee-clone/shopee/services/auth/internal/infrastructure/hash"
	"github.com/shopee-clone/shopee/services/auth/internal/infrastructure/jwt"
	"github.com/shopee-clone/shopee/services/auth/internal/infrastructure/mysql"
	redisinfra "github.com/shopee-clone/shopee/services/auth/internal/infrastructure/redis"
	"github.com/shopee-clone/shopee/services/auth/internal/security"
	"github.com/shopee-clone/shopee/services/auth/internal/tracing"
	grpctransport "github.com/shopee-clone/shopee/services/auth/internal/transport/grpc"
	httptransport "github.com/shopee-clone/shopee/services/auth/internal/transport/http"
	pb "github.com/shopee-clone/shopee/services/auth/proto/auth/v1"
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
		logger.Warn("redis not available, continuing without redis", zap.Error(err))
		redisClient = nil
	}

	userRepo := mysql.NewUserRepository(db)
	sessionRepo := mysql.NewSessionRepository(db)
	auditRepo := mysql.NewAuditRepository(db, cfg.Audit)
	defer auditRepo.Stop()

	redisStore := redisinfra.NewStore(redisClient, cfg.Redis)

	hashService := hash.NewService(cfg.Password)

	jwtService := jwt.NewService(cfg.JWT, redisStore)
	if redisClient == nil {
		jwtService = jwt.NewService(cfg.JWT, nil)
	}

	rateLimiter := security.NewRateLimiter(redisClient, cfg.RateLimit)
	suspiciousDetector := security.NewSuspiciousDetector(redisClient, cfg.Security)

	authService := application.NewAuthService(
		cfg, userRepo, sessionRepo, auditRepo, redisStore,
		rateLimiter, suspiciousDetector, jwtService, hashService,
	)

	gin.SetMode(getGinMode(cfg.AppEnv))
	engine := gin.New()

	handler := httptransport.NewHandler(authService)

	httpRouter := httptransport.NewRouter(handler, nil)
	httpRouter.Setup(engine)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	grpcServer := grpc.NewServer()
	grpcAuthServer := grpctransport.NewAuthGRPCServer(authService)
	pb.RegisterAuthServiceServer(grpcServer, grpcAuthServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, &grpcHealthServer{})
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen grpc", zap.Int("port", cfg.GRPCPort), zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("starting auth service",
			zap.Int("http_port", cfg.HTTPPort),
			zap.Int("grpc_port", cfg.GRPCPort),
			zap.String("env", cfg.AppEnv),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("grpc server failed", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("shutting down auth service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown error", zap.Error(err))
	}

	if redisClient != nil {
		redisClient.Close()
	}

	logger.Info("auth service stopped")
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
	case "production", "staging":
		return gin.ReleaseMode
	default:
		return gin.DebugMode
	}
}
