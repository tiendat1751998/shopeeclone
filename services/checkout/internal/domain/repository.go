package domain

import (
	"context"
	"time"
)

type CheckoutRepository interface {
	FindByID(ctx context.Context, id string) (*Checkout, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*Checkout, error)
	Create(ctx context.Context, c *Checkout) error
	Update(ctx context.Context, c *Checkout) error
	UpdateStatus(ctx context.Context, id, status string) error
	FindExpired(ctx context.Context, before time.Time, limit int) ([]*Checkout, error)
}

type CheckoutStepLogRepository interface {
	Create(ctx context.Context, log *CheckoutStepLog) error
	FindByCheckoutID(ctx context.Context, checkoutID string) ([]*CheckoutStepLog, error)
}

type PricingSnapshotRepository interface {
	FindByCheckoutID(ctx context.Context, checkoutID string) (*PricingSnapshot, error)
	Create(ctx context.Context, snap *PricingSnapshot) error
}

type ReservationOrchestrationRepository interface {
	FindByCheckoutID(ctx context.Context, checkoutID string) ([]*ReservationOrchestration, error)
	Create(ctx context.Context, r *ReservationOrchestration) error
	UpdateStatus(ctx context.Context, id, status string) error
}

type ReconciliationJobRepository interface {
	Create(ctx context.Context, job *ReconciliationJob) error
	FindPending(ctx context.Context, limit int) ([]*ReconciliationJob, error)
	Update(ctx context.Context, job *ReconciliationJob) error
	IncrementAttempt(ctx context.Context, id string) error
}
