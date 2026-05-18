package domain

import "context"

type VoucherRepository interface {
	FindByID(ctx context.Context, id string) (*Voucher, error)
	FindByCode(ctx context.Context, code string) (*Voucher, error)
	Create(ctx context.Context, v *Voucher) error
	Update(ctx context.Context, v *Voucher) error
	IncrementUsage(ctx context.Context, id string) error
	ListActive(ctx context.Context, offset, limit int) ([]*Voucher, int64, error)
}

type VoucherRedemptionRepository interface {
	FindByID(ctx context.Context, id string) (*VoucherRedemption, error)
	FindByUserAndVoucher(ctx context.Context, userID, voucherID string) ([]*VoucherRedemption, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*VoucherRedemption, error)
	CountByVoucher(ctx context.Context, voucherID string) (int64, error)
	CountByUserAndVoucher(ctx context.Context, userID, voucherID string) (int, error)
	Create(ctx context.Context, r *VoucherRedemption) error
}

type CampaignRepository interface {
	FindByID(ctx context.Context, id string) (*Campaign, error)
	Create(ctx context.Context, c *Campaign) error
	Update(ctx context.Context, c *Campaign) error
	ListActive(ctx context.Context) ([]*Campaign, error)
	ListByType(ctx context.Context, cType string, offset, limit int) ([]*Campaign, int64, error)
}

type PricingRuleRepository interface {
	FindByCampaign(ctx context.Context, campaignID string) ([]*PricingRule, error)
	Create(ctx context.Context, rule *PricingRule) error
	Update(ctx context.Context, rule *PricingRule) error
}

type EligibilityRuleRepository interface {
	FindByPromotion(ctx context.Context, promotionID string) ([]*EligibilityRule, error)
	Create(ctx context.Context, rule *EligibilityRule) error
}

type StackingRuleRepository interface {
	FindByPromotionType(ctx context.Context, pType string) ([]*StackingRule, error)
	Create(ctx context.Context, rule *StackingRule) error
}
