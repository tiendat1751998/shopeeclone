package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/platforms/payment-ledger/internal/dispute"
	"github.com/shopee-clone/shopee/platforms/payment-ledger/internal/ledger"
	"github.com/shopee-clone/shopee/platforms/payment-ledger/internal/payment"
	"github.com/shopee-clone/shopee/platforms/payment-ledger/internal/payout"
	"github.com/shopee-clone/shopee/platforms/payment-ledger/internal/reconciliation"
	httptransport "github.com/shopee-clone/shopee/platforms/payment-ledger/internal/transport/http"
)

func main() {
	paymentRepo := payment.NewInMemoryRepository()
	paymentService := payment.NewService(paymentRepo)

	accountRepo := ledger.NewInMemoryAccountRepo()
	txnRepo := ledger.NewInMemoryTransactionRepo()
	entryRepo := ledger.NewInMemoryEntryRepo()
	ledgerService := ledger.NewService(accountRepo, txnRepo, entryRepo)

	payoutRepo := payout.NewInMemoryPayoutRepo()
	batchRepo := payout.NewInMemoryBatchRepo()
	payoutService := payout.NewService(payoutRepo, batchRepo)

	reconRepo := reconciliation.NewInMemoryRepository()
	reconService := reconciliation.NewService(reconRepo, nil, nil)

	disputeRepo := dispute.NewInMemoryRepository()
	disputeService := dispute.NewService(disputeRepo)

	handler := httptransport.NewHandler(paymentService, ledgerService, payoutService, reconService, disputeService)
	router := httptransport.NewRouter(handler)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	router.Setup(engine)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("starting payment-ledger server on :8080")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
