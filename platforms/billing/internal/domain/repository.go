package domain

import "context"

type AccountRepository interface {
	Create(ctx context.Context, a *Account) error
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByUserAndCurrency(ctx context.Context, userID string, currency string) (*Account, error)
	UpdateBalance(ctx context.Context, id string, balance, frozen int64) error
}

type TransactionRepository interface {
	Create(ctx context.Context, txn *Transaction) error
	GetByID(ctx context.Context, id string) (*Transaction, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*Transaction, error)
	UpdateStatus(ctx context.Context, id string, status TransactionStatus) error
	ListByAccount(ctx context.Context, accountID string, offset, limit int) ([]*Transaction, int64, error)
}

type LedgerEntryRepository interface {
	CreateBatch(ctx context.Context, entries []*LedgerEntry) error
	GetByTransaction(ctx context.Context, txnID string) ([]*LedgerEntry, error)
	GetByAccount(ctx context.Context, accountID string, offset, limit int) ([]*LedgerEntry, int64, error)
}

type WalletRepository interface {
	Create(ctx context.Context, w *Wallet) error
	GetByID(ctx context.Context, id string) (*Wallet, error)
	GetByUser(ctx context.Context, userID string) ([]*Wallet, error)
	GetByUserAndType(ctx context.Context, userID, walletType string) (*Wallet, error)
	UpdateBalance(ctx context.Context, id string, balance, frozen, pending int64) error
}

type SettlementRepository interface {
	Create(ctx context.Context, s *Settlement) error
	GetByID(ctx context.Context, id string) (*Settlement, error)
	UpdateStatus(ctx context.Context, id string, status SettlementStatus) error
	ListPending(ctx context.Context, limit int) ([]*Settlement, error)
}

type PayoutRepository interface {
	Create(ctx context.Context, p *Payout) error
	GetByID(ctx context.Context, id string) (*Payout, error)
	UpdateStatus(ctx context.Context, id string, status PayoutStatus) error
	ListPending(ctx context.Context, limit int) ([]*Payout, error)
}

type RefundRepository interface {
	Create(ctx context.Context, r *Refund) error
	GetByID(ctx context.Context, id string) (*Refund, error)
	GetByOriginalTxn(ctx context.Context, txnID string) (*Refund, error)
	UpdateStatus(ctx context.Context, id string, status RefundStatus) error
}

type ReconciliationRepository interface {
	CreateRun(ctx context.Context, r *ReconciliationRun) error
	GetRun(ctx context.Context, id string) (*ReconciliationRun, error)
}

type AuditLogRepository interface {
	Append(ctx context.Context, log *AuditLog) error
	ListByResource(ctx context.Context, resource, resourceID string, offset, limit int) ([]*AuditLog, int64, error)
}
