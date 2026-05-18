package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Checkout represents the checkout orchestration state
type Checkout struct {
	ID              string    `db:"id" json:"id"`
	UserID          string    `db:"user_id" json:"user_id"`
	CartID          string    `db:"cart_id" json:"cart_id"`
	OrderID         string    `db:"order_id" json:"order_id,omitempty"`
	Status          string    `db:"status" json:"status"`
	IdempotencyKey  string    `db:"idempotency_key" json:"idempotency_key"`
	CurrentStep     string    `db:"current_step" json:"current_step"`
	FailureReason   string    `db:"failure_reason" json:"failure_reason,omitempty"`
	AttemptCount    int       `db:"attempt_count" json:"attempt_count"`
	ReservationKeys string   `db:"reservation_keys" json:"reservation_keys"`
	PricingSnapshot string    `db:"pricing_snapshot" json:"pricing_snapshot"`
	PromotionResults string   `db:"promotion_results" json:"promotion_results"`
	Subtotal        int64     `db:"subtotal" json:"subtotal"`
	DiscountTotal   int64     `db:"discount_total" json:"discount_total"`
	ShippingTotal   int64     `db:"shipping_total" json:"shipping_total"`
	GrandTotal      int64     `db:"grand_total" json:"grand_total"`
	Currency        string    `db:"currency" json:"currency"`
	ExpiresAt       time.Time `db:"expires_at" json:"expires_at"`
	CompletedAt     *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// Checkout states
const (
	CheckoutStatusPending       = "pending"
	CheckoutStatusValidating    = "validating"
	CheckoutStatusPricingFrozen = "pricing_frozen"
	CheckoutStatusReserving     = "reserving_inventory"
	CheckoutStatusReserved      = "inventory_reserved"
	CheckoutStatusProcessing    = "processing_payment"
	CheckoutStatusCompleted     = "completed"
	CheckoutStatusFailed        = "failed"
	CheckoutStatusRollingBack   = "rolling_back"
	CheckoutStatusRolledBack    = "rolled_back"
	CheckoutStatusExpired       = "expired"
)

// Checkout steps
const (
	StepInit         = "init"
	StepValidate     = "validate"
	StepFreezePricing = "freeze_pricing"
	StepReserve      = "reserve_inventory"
	StepProcess      = "process_payment"
	StepComplete     = "complete"
	StepRollback     = "rollback"
)

func NewCheckout(userID, cartID, idempotencyKey string, ttl time.Duration) *Checkout {
	now := time.Now()
	return &Checkout{
		ID:             uuid.New().String(),
		UserID:         userID,
		CartID:         cartID,
		Status:         CheckoutStatusPending,
		IdempotencyKey: idempotencyKey,
		CurrentStep:    StepInit,
		AttemptCount:   0,
		ExpiresAt:      now.Add(ttl),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (c *Checkout) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

func (c *Checkout) CanRetry() bool {
	return c.AttemptCount < 3 && c.Status != CheckoutStatusCompleted && c.Status != CheckoutStatusRolledBack
}

func (c *Checkout) AdvanceStep(step string) {
	c.CurrentStep = step
	c.UpdatedAt = time.Now()
}

func (c *Checkout) IncrementAttempt() {
	c.AttemptCount++
	c.UpdatedAt = time.Now()
}

func (c *Checkout) MarkCompleted(orderID string) {
	c.OrderID = orderID
	c.Status = CheckoutStatusCompleted
	c.CurrentStep = StepComplete
	now := time.Now()
	c.CompletedAt = &now
	c.UpdatedAt = now
}

func (c *Checkout) MarkFailed(reason string) {
	c.Status = CheckoutStatusFailed
	c.FailureReason = reason
	c.UpdatedAt = time.Now()
}

func (c *Checkout) MarkRollingBack() {
	c.Status = CheckoutStatusRollingBack
	c.CurrentStep = StepRollback
	c.UpdatedAt = time.Now()
}

func (c *Checkout) MarkRolledBack() {
	c.Status = CheckoutStatusRolledBack
	c.UpdatedAt = time.Now()
}

// CheckoutStepLog records each step execution for audit/debugging
type CheckoutStepLog struct {
	ID         string    `db:"id" json:"id"`
	CheckoutID string    `db:"checkout_id" json:"checkout_id"`
	Step       string    `db:"step" json:"step"`
	Status     string    `db:"status" json:"status"`
	ErrorMsg   string    `db:"error_message" json:"error_message,omitempty"`
	Metadata   string    `db:"metadata" json:"metadata,omitempty"`
	DurationMs int64     `db:"duration_ms" json:"duration_ms"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

func NewCheckoutStepLog(checkoutID, step, status string, durationMs int64, errMsg, metadata string) *CheckoutStepLog {
	return &CheckoutStepLog{
		ID:         uuid.New().String(),
		CheckoutID: checkoutID,
		Step:       step,
		Status:     status,
		ErrorMsg:   errMsg,
		Metadata:   metadata,
		DurationMs: durationMs,
		CreatedAt:  time.Now(),
	}
}

// PricingSnapshot stores frozen pricing at checkout time
type PricingSnapshot struct {
	ID              string    `db:"id" json:"id"`
	CheckoutID      string    `db:"checkout_id" json:"checkout_id"`
	Items           string    `db:"items" json:"items"`
	SellerGroups    string    `db:"seller_groups" json:"seller_groups"`
	Subtotal        int64     `db:"subtotal" json:"subtotal"`
	DiscountTotal   int64     `db:"discount_total" json:"discount_total"`
	ShippingTotal   int64     `db:"shipping_total" json:"shipping_total"`
	GrandTotal      int64     `db:"grand_total" json:"grand_total"`
	Currency        string    `db:"currency" json:"currency"`
	PromotionsApplied string  `db:"promotions_applied" json:"promotions_applied"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// ReservationOrchestration tracks inventory reservations for a checkout
type ReservationOrchestration struct {
	ID              string    `db:"id" json:"id"`
	CheckoutID      string    `db:"checkout_id" json:"checkout_id"`
	ReservationKey  string    `db:"reservation_key" json:"reservation_key"`
	SKU             string    `db:"sku" json:"sku"`
	WarehouseID     string    `db:"warehouse_id" json:"warehouse_id"`
	Quantity        int64     `db:"quantity" json:"quantity"`
	Status          string    `db:"status" json:"status"`
	ErrorMessage    string    `db:"error_message" json:"error_message,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

const (
	ReservationStatusPending   = "pending"
	ReservationStatusReserved = "reserved"
	ReservationStatusReleased = "released"
	ReservationStatusFailed   = "failed"
)

// ReconciliationJob tracks async reconciliation tasks
type ReconciliationJob struct {
	ID           string    `db:"id" json:"id"`
	CheckoutID   string    `db:"checkout_id" json:"checkout_id"`
	JobType      string    `db:"job_type" json:"job_type"`
	Status       string    `db:"status" json:"status"`
	AttemptCount int       `db:"attempt_count" json:"attempt_count"`
	MaxAttempts  int       `db:"max_attempts" json:"max_attempts"`
	NextRetryAt  time.Time `db:"next_retry_at" json:"next_retry_at"`
	Metadata     string    `db:"metadata" json:"metadata"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

const (
	JobTypeReleaseReservation = "release_reservation"
	JobTypeConfirmReservation = "confirm_reservation"
	JobTypeUpdateOrderStatus  = "update_order_status"

	JobStatusPending  = "pending"
	JobStatusRunning  = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed   = "failed"
)

// Domain errors
var (
	ErrCheckoutNotFound     = ErrCheckout("checkout_not_found")
	ErrCheckoutExpired      = ErrCheckout("checkout_expired")
	ErrCheckoutCompleted    = ErrCheckout("checkout_already_completed")
	ErrIdempotencyConflict  = ErrCheckout("idempotency_conflict")
	ErrValidationFailed     = ErrCheckout("validation_failed")
	ErrPricingChanged       = ErrCheckout("pricing_changed")
	ErrReservationFailed    = ErrCheckout("reservation_failed")
	ErrRollbackFailed       = ErrCheckout("rollback_failed")
	ErrMaxRetriesExceeded   = ErrCheckout("max_retries_exceeded")
)

type ErrCheckout string

func (e ErrCheckout) Error() string {
	return "checkout: " + string(e)
}
