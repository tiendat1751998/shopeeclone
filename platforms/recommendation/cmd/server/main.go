package main
import ("context"; "fmt"; "net/http"; "os"; "os/signal"; "syscall"; "time"; "github.com/gin-gonic/gin"; sharedRedis "github.com/shopee-clone/shopee/packages/go-shared/pkg/redis"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/health"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "github.com/shopee-clone/shopee/platforms/recommendation/internal/application"; "github.com/shopee-clone/shopee/platforms/recommendation/internal/config"; redisinfra "github.com/shopee-clone/shopee/platforms/recommendation/internal/infrastructure/redis"; httptransport "github.com/shopee-clone/shopee/platforms/recommendation/internal/transport/http"; kafkatransport "github.com/shopee-clone/shopee/platforms/recommendation/internal/transport/kafka"; "go.uber.org/zap")
var version = "1.0.0"
func main() {
	cfg := config.Load(); logger := observability.InitLogger(cfg.AppName, cfg.LogLevel)
	shutdownTracer, _ := observability.InitTracer(cfg.OpenTelemetry.ServiceName, cfg.OpenTelemetry.Endpoint)
	defer shutdownTracer(); defer observability.Sync()
	redisClient, err := sharedRedis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil { logger.Warn("redis not available", zap.Error(err)); redisClient = nil }
	redisStore := redisinfra.NewStore(redisClient, cfg.Redis)
	var pub application.EventPublisher
	if len(cfg.Kafka.Brokers) > 0 && cfg.Kafka.Brokers[0] != "" { p := kafkatransport.NewProducer(cfg.Kafka.Brokers, cfg.AppName); defer p.Close(); pub = p }
	recService := application.NewRecommendationService(nil, nil, nil, redisStore, pub)
	hc := health.NewChecker(cfg.AppName, version, redisClient)
	gin.SetMode(getGinMode(cfg.AppEnv)); engine := gin.New()
	handler := httptransport.NewHandler(recService); router := httptransport.NewRouter(handler, hc); router.Setup(engine)
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", cfg.HTTPPort), Handler: engine, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	quit := make(chan os.Signal, 1); signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() { logger.Info("starting recommendation service", zap.Int("http_port", cfg.HTTPPort)); if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed { logger.Fatal("http server failed", zap.Error(err)) } }()
	<-quit; logger.Info("shutting down..."); ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second); defer cancel()
	httpServer.Shutdown(ctx); if redisClient != nil { redisClient.Close() }; logger.Info("stopped")
}
func getGinMode(env string) string { if env == "production" || env == "staging" { return gin.ReleaseMode }; return gin.DebugMode }
