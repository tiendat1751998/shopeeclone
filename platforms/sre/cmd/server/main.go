package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"go.uber.org/zap"

	"github.com/shopee-clone/shopee/platforms/sre/internal/alerting"
	"github.com/shopee-clone/shopee/platforms/sre/internal/deployment"
	"github.com/shopee-clone/shopee/platforms/sre/internal/healthcheck"
	"github.com/shopee-clone/shopee/platforms/sre/internal/incident"
	"github.com/shopee-clone/shopee/platforms/sre/internal/runbook"
	"github.com/shopee-clone/shopee/platforms/sre/internal/slo"
	httptransport "github.com/shopee-clone/shopee/platforms/sre/internal/transport/http"
)

var version = "1.0.0"
var port = 8080

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	incidentRepo := incident.NewInMemoryRepository()
	incidentSvc := incident.NewService(incidentRepo)

	alertingRepo := alerting.NewInMemoryRepository()
	alertingSvc := alerting.NewService(alertingRepo)

	healthcheckRepo := healthcheck.NewInMemoryRepository()
	healthcheckSvc := healthcheck.NewService(healthcheckRepo)

	sloRepo := slo.NewInMemoryRepository()
	sloSvc := slo.NewService(sloRepo)

	deploymentRepo := deployment.NewInMemoryRepository()
	deploymentSvc := deployment.NewService(deploymentRepo)

	runbookRepo := runbook.NewInMemoryRepository()
	runbookSvc := runbook.NewService(runbookRepo)

	hc := health.NewChecker("shopee-sre", version)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	handler := httptransport.NewHandler(incidentSvc, alertingSvc, healthcheckSvc, sloSvc, deploymentSvc, runbookSvc)
	router := httptransport.NewRouter(handler, hc)
	router.Setup(engine)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting sre platform", zap.Int("http_port", port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	httpServer.Shutdown(shutdownCtx)
	wg.Wait()
	logger.Info("stopped")
}
