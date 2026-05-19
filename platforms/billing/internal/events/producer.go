package events

import "context"

const (
	EventLedgerPosted       = "ledger.posted"
	EventTransactionCreated = "transaction.created"
	EventTransactionCompleted = "transaction.completed"
	EventTransactionFailed  = "transaction.failed"
	EventWalletUpdated      = "wallet.updated"
	EventSettlementCreated  = "settlement.created"
	EventSettlementCompleted = "settlement.completed"
	EventPayoutInitiated    = "payout.initiated"
	EventPayoutCompleted    = "payout.completed"
	EventRefundProcessed    = "refund.processed"
	EventReconciliationMismatch = "reconciliation.mismatch"
)

type Publisher interface {
	Publish(ctx context.Context, eventType string, payload interface{}) error
}
