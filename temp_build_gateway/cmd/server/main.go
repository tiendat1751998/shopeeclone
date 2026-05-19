package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	sharedRedis "github.com/shopee-clone/shopee/packages/go-shared/pkg/redis"
	"github.com/shopee-clone/shopee/services/gateway/internal/auth"
	"github.com/shopee-clone/shopee/services/gateway/internal/config"
	"github.com/shopee-clone/shopee/services/gateway/internal/discovery"
	"github.com/shopee-clone/shopee/services/gateway/internal/ratelimit"
	"github.com/shopee-clone/shopee/services/gateway/internal/routing"
	"github.com/shopee-clone/shopee/services/gateway/internal/tracing"
	"github.com/shopee-clone/shopee/services/gateway/internal/transport"
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

	redisClient, err := sharedRedis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Warn("redis not available, continuing without redis", zap.Error(err))
		redisClient = nil
	}

	jwtValidator := auth.NewJWTValidator(cfg.Auth, redisClient)
	authMiddleware := auth.NewAuthMiddleware(jwtValidator)

	rateLimiter := ratelimit.NewRateLimiter(redisClient, cfg.RateLimit)

	svcDiscovery := discovery.NewServiceDiscovery()
	registerUpstreams(cfg, svcDiscovery)

	proxy := transport.NewProxy(
		svcDiscovery,
		cfg.Upstreams.MaxIdleConns,
		cfg.Upstreams.IdleConnTimeout,
	)

	healthChecker := health.NewChecker(cfg.AppName, version)
	registerHealthChecks(healthChecker, redisClient)

	gin.SetMode(getGinMode(cfg.AppEnv))
	engine := gin.New()

	router := routing.NewRouter(cfg, proxy, rateLimiter, authMiddleware, svcDiscovery, healthChecker)
	router.Setup(engine)

	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:        engine,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelUnaryServerInterceptor()),
	)
	grpc_health_v1.RegisterHealthServer(grpcServer, &gatewayHealthServer{})
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen gRPC", zap.Int("port", cfg.GRPCPort), zap.Error(err))
	}

	go func() {
		logger.Info("starting gateway",
			zap.Int("http_port", cfg.HTTPPort),
			zap.Int("grpc_port", cfg.GRPCPort),
			zap.String("env", cfg.AppEnv),
			zap.String("version", version),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("gateway http server failed", zap.Error(err))
		}
	}()

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("gateway grpc server failed", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("shutting down gateway...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("gateway http shutdown error", zap.Error(err))
	}

	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			logger.Error("redis close error", zap.Error(err))
		}
	}

	logger.Info("gateway stopped")
}

func registerUpstreams(cfg *config.Config, svcDiscovery *discovery.ServiceDiscovery) {
	upstreams := map[string]*discovery.ServiceInstance{
		"auth":           {ID: "auth-1", Name: "auth", Address: cfg.Upstreams.AuthService, Port: 8080, Weight: 10},
		"catalog":        {ID: "catalog-1", Name: "catalog", Address: cfg.Upstreams.CatalogService, Port: 8080, Weight: 10},
		"cart":           {ID: "cart-1", Name: "cart", Address: cfg.Upstreams.CartService, Port: 8080, Weight: 10},
		"order":          {ID: "order-1", Name: "order", Address: cfg.Upstreams.OrderService, Port: 8080, Weight: 10},
		"inventory":      {ID: "inventory-1", Name: "inventory", Address: cfg.Upstreams.InventoryService, Port: 8080, Weight: 10},
		"payment":        {ID: "payment-1", Name: "payment", Address: cfg.Upstreams.PaymentService, Port: 8080, Weight: 10},
		"search":         {ID: "search-1", Name: "search", Address: cfg.Upstreams.SearchService, Port: 8080, Weight: 10},
		"recommendation": {ID: "rec-1", Name: "recommendation", Address: cfg.Upstreams.RecommendationService, Port: 8080, Weight: 5},
	}

	for name, instance := range upstreams {
		svcDiscovery.RegisterStatic(name, []*discovery.ServiceInstance{instance})
	}
}

func registerHealthChecks(healthChecker *health.Checker, redisClient *redis.Client) {
	if redisClient != nil {
		healthChecker.AddCheck("redis", func(ctx context.Context) error {
			return redisClient.Ping(ctx).Err()
		})
	}
}

func otelUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, span := tracing.Tracer.Start(ctx, info.FullMethod)
		defer span.End()
		return handler(ctx, req)
	}
}

type gatewayHealthServer struct{}

func (s *gatewayHealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *gatewayHealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func getGinMode(env string) string {
	switch env {
	case "production":
		return gin.ReleaseMode
	case "staging":
		return gin.ReleaseMode
	default:
		return gin.DebugMode
	}
}
