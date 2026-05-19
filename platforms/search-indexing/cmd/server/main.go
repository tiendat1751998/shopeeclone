package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/bulkindexer"
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/coordinator"
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/monitoring"
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/pipeline"
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/synonyms"
	httptransport "github.com/shopee-clone/shopee/platforms/search-indexing/internal/transport/http"
)

func main() {
	coordRepo := coordinator.NewInMemoryRepository()
	coordSvc := coordinator.NewService(coordRepo)

	bulkRepo := bulkindexer.NewInMemoryRepository()
	bulkSvc := bulkindexer.NewService(bulkRepo)

	pipelineRepo := pipeline.NewInMemoryRepository()
	pipelineSvc := pipeline.NewService(pipelineRepo)

	synonymRepo := synonyms.NewInMemoryRepository()
	synonymSvc := synonyms.NewService(synonymRepo)

	monitorRepo := monitoring.NewInMemoryRepository()
	monitorSvc := monitoring.NewService(monitorRepo)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	handler := httptransport.NewHandler(coordSvc, bulkSvc, pipelineSvc, synonymSvc, monitorSvc)
	router := httptransport.NewRouter(handler)
	router.Setup(engine)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("starting search-indexing service on :%s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	httpServer.Shutdown(shutdownCtx)
	wg.Wait()
	log.Println("stopped")
}
