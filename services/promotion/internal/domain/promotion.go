package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Voucher represents a discount voucher
type Voucher struct {
	ID             string    `db:"id" json:"id"`
	Code           string    `db:"code" json:"code"`
	Title          string    `db:"title" json:"title"`
	Description    string    `db:"description" json:"description,omitempty"`
	Type           string    `db:"type" json:"type"`
	DiscountValue  int64     `db:"discount_value" json:"discount_value"`
	MinSpend       int64     `db:"min_spend" json:"min_spend"`
	MaxDiscount    int64     `db:"max_discount" json:"max_discount"`
	UsageLimit     int64     `db:"usage_limit" json:"usage_limit"`
	UsageCount     int64     `db:"usage_count" json:"usage_count"`
	PerUserLimit   int       `db:"per_user_limit" json:"per_user_limit"`
	Scope          string    `db:"scope" json:"scope"`
	ShopID         *string   `db:"shop_id" json:"shop_id,omitempty"`
	CategoryID     *string   `db:"category_id" json:"category_id,omitempty"`
	SKU            *string   `db:"sku" json:"sku,omitempty"`
	Region         *string   `db:"region" json:"region,omitempty"`
	PaymentMethod  *string   `db:"payment_method" json:"payment_method,omitempty"`
	StartTime      time.Time `db:"start_time" json:"start_time"`
	EndTime        time.Time `db:"end_time" json:"end_time"`
	Status         string    `db:"status" json:"status"`
	Stackable      bool      `db:"stackable" json:"stackable"`
	Priority       int       `db:"priority" json:"priority"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

const (
	VoucherTypePercentage = "percentage"
	VoucherTypeFixed      = "fixed"
	VoucherTypeShipping   = "shipping"

	VoucherScopePlatform  = "platform"
	VoucherScopeShop      = "shop"
	VoucherScopeCategory  = "category"
	VoucherScopeSKU       = "sku"

	VoucherStatusActive   = "active"
	VoucherStatusInactive = "inactive"
	VoucherStatusExpired  = "expired"
	VoucherStatusExhausted = "exhausted"
)

func NewVoucher(code, title, vType string, discountValue, minSpend, maxDiscount int64, startTime, endTime time.Time) *Voucher {
	now := time.Now()
	return &Voucher{
		ID:            uuid.New().String(),
		Code:          code,
		Title:         title,
		Type:          vType,
		DiscountValue: discountValue,
		MinSpend:      minSpend,
		MaxDiscount:   maxDiscount,
		UsageLimit:    10000,
		PerUserLimit:  1,
		Scope:         VoucherScopePlatform,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        VoucherStatusActive,
		Stackable:     false,
		Priority:      0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (v *Voucher) IsActive() bool {
	now := time.Now()
	return v.Status == VoucherStatusActive && now.After(v.StartTime) && now.Before(v.EndTime) && v.UsageCount < v.UsageLimit
}

func (v *Voucher) IsExpired() bool {
	return time.Now().After(v.EndTime)
}

func (v *Voucher) CanRedeem(userID string, subtotal int64, shopID, categoryID, sku, region, paymentMethod string) error {
	if !v.IsActive() {
		return fmt.Errorf("%w: voucher not active", ErrVoucherInvalid)
	}
	if subtotal < v.MinSpend {
		return fmt.Errorf("%w: subtotal %d < min spend %d", ErrVoucherMinSpend, subtotal, v.MinSpend)
	}
	if v.Scope == VoucherScopeShop && (v.ShopID == nil || *v.ShopID != shopID) {
		return fmt.Errorf("%w: voucher not valid for this shop", ErrVoucherScope)
	}
	if v.Scope == VoucherScopeCategory && (v.CategoryID == nil || *v.CategoryID != categoryID) {
		return fmt.Errorf("%w: voucher not valid for this category", ErrVoucherScope)
	}
	if v.Scope == VoucherScopeSKU && (v.SKU == nil || *v.SKU != sku) {
		return fmt.Errorf("%w: voucher not valid for this SKU", ErrVoucherScope)
	}
	if v.Region != nil && *v.Region != "" && *v.Region != region {
		return fmt.Errorf("%w: voucher not valid for this region", ErrVoucherRegion)
	}
	if v.PaymentMethod != nil && *v.PaymentMethod != "" && *v.PaymentMethod != paymentMethod {
		return fmt.Errorf("%w: voucher not valid for this payment method", ErrVoucherPaymentMethod)
	}
	return nil
}

func (v *Voucher) CalculateDiscount(subtotal int64) int64 {
	var discount int64
	switch v.Type {
	case VoucherTypePercentage:
		discount = subtotal * v.DiscountValue / 100
		if discount > v.MaxDiscount && v.MaxDiscount > 0 {
			discount = v.MaxDiscount
		}
	case VoucherTypeFixed:
		discount = v.DiscountValue
		if discount > subtotal {
			discount = subtotal
		}
	case VoucherTypeShipping:
		discount = v.DiscountValue
	}
	return discount
}

// VoucherRedemption tracks voucher usage
type VoucherRedemption struct {
	ID           string    `db:"id" json:"id"`
	VoucherID    string    `db:"voucher_id" json:"voucher_id"`
	UserID       string    `db:"user_id" json:"user_id"`
	OrderID      string    `db:"order_id" json:"order_id"`
	DiscountAmt  int64     `db:"discount_amount" json:"discount_amount"`
	IdempotencyKey string  `db:"idempotency_key" json:"idempotency_key"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// Campaign represents a promotional campaign
type Campaign struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Type        string    `db:"type" json:"type"`
	Description string    `db:"description" json:"description,omitempty"`
	Rules       string    `db:"rules" json:"rules"`
	StartTime   time.Time `db:"start_time" json:"start_time"`
	EndTime     time.Time `db:"end_time" json:"end_time"`
	Status      string    `db:"status" json:"status"`
	Priority    int       `db:"priority" json:"priority"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

const (
	CampaignTypeFlashSale    = "flash_sale"
	CampaignTypeScheduled    = "scheduled"
	CampaignTypeSeasonal     = "seasonal"
	CampaignTypeSeller       = "seller"
	CampaignTypeCategory     = "category"

	CampaignStatusDraft     = "draft"
	CampaignStatusActive    = "active"
	CampaignStatusPaused    = "paused"
	CampaignStatusEnded     = "ended"
)

// PricingRule defines discount calculation rules
type PricingRule struct {
	ID            string    `db:"id" json:"id"`
	CampaignID    string    `db:"campaign_id" json:"campaign_id"`
	RuleType      string    `db:"rule_type" json:"rule_type"`
	Condition     string    `db:"condition" json:"condition"`
	Action        string    `db:"action" json:"action"`
	Priority      int       `db:"priority" json:"priority"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// EligibilityRule defines who can use a promotion
type EligibilityRule struct {
	ID           string `db:"id" json:"id"`
	PromotionID  string `db:"promotion_id" json:"promotion_id"`
	TargetType   string `db:"target_type" json:"target_type"`
	TargetValue  string `db:"target_value" json:"target_value"`
	IsActive     bool   `db:"is_active" json:"is_active"`
}

// StackingRule defines how promotions can be combined
type StackingRule struct {
	ID              string `db:"id" json:"id"`
	PromotionType   string `db:"promotion_type" json:"promotion_type"`
	CanStackWith    string `db:"can_stack_with" json:"can_stack_with"`
	MaxStackCount   int    `db:"max_stack_count" json:"max_stack_count"`
	Priority        int    `db:"priority" json:"priority"`
}

// PromotionResult represents the result of promotion evaluation
type PromotionResult struct {
	VoucherID    string `json:"voucher_id,omitempty"`
	CampaignID   string `json:"campaign_id,omitempty"`
	Type         string `json:"type"`
	Description  string `json:"description"`
	DiscountAmt  int64  `json:"discount_amount"`
	FinalAmount  int64  `json:"final_amount"`
}

// Domain errors
var (
	ErrVoucherInvalid       = ErrPromotion("voucher_invalid")
	ErrVoucherExpired       = ErrPromotion("voucher_expired")
	ErrVoucherExhausted     = ErrPromotion("voucher_exhausted")
	ErrVoucherMinSpend      = ErrPromotion("voucher_min_spend_not_met")
	ErrVoucherScope         = ErrPromotion("voucher_scope_mismatch")
	ErrVoucherRegion        = ErrPromotion("voucher_region_mismatch")
	ErrVoucherPaymentMethod = ErrPromotion("voucher_payment_method_mismatch")
	ErrVoucherUserLimit     = ErrPromotion("voucher_user_limit_reached")
	ErrDuplicateRedemption  = ErrPromotion("duplicate_redemption")
	ErrCampaignNotFound     = ErrPromotion("campaign_not_found")
	ErrPromotionNotFound    = ErrPromotion("promotion_not_found")
	ErrStackingConflict     = ErrPromotion("stacking_conflict")
)

type ErrPromotion string

func (e ErrPromotion) Error() string {
	return "promotion: " + string(e)
}
