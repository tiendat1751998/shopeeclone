package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/shopee-clone/shopee/platforms/service-mesh/internal/discovery"
	"github.com/shopee-clone/shopee/platforms/service-mesh/internal/loadbalancer"
	"github.com/shopee-clone/shopee/platforms/service-mesh/internal/mtls"
	"github.com/shopee-clone/shopee/platforms/service-mesh/internal/resilience"
	"github.com/shopee-clone/shopee/platforms/service-mesh/internal/telemetry"
	"github.com/shopee-clone/shopee/platforms/service-mesh/internal/traffic"
	httpTransport "github.com/shopee-clone/shopee/platforms/service-mesh/internal/transport/http"
)

func main() {
	discRepo := discovery.NewInMemoryRepository()
	discSvc := discovery.NewService(discRepo)

	ca, err := mtls.NewCertificateAuthority("ShopeeClone", "Service Mesh CA", 3650)
	if err != nil {
		panic(err)
	}
	certMgr := mtls.NewCertManager(ca)

	trafficRepo := traffic.NewInMemoryRepository()
	trafficEng := traffic.NewEngine(trafficRepo)

	lb := loadbalancer.NewLoadBalancer(loadbalancer.RoundRobin)

	executor := resilience.NewExecutor()

	telRepo := telemetry.NewInMemoryRepository()

	handler := httpTransport.NewHandler(discSvc, certMgr, trafficEng, lb, executor, telRepo)
	router := httpTransport.NewRouter(handler)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	router.Setup(engine)

	srv := &http.Server{
		Addr:         ":8095",
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
