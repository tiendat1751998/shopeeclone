package mysql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/shopee-clone/shopee/services/promotion/internal/domain"
)

type VoucherRepository struct{ db *sqlx.DB }
func NewVoucherRepository(db *sqlx.DB) *VoucherRepository { return &VoucherRepository{db: db} }

func (r *VoucherRepository) FindByID(ctx context.Context, id string) (*domain.Voucher, error) {
	var v domain.Voucher
	err := r.db.GetContext(ctx, &v, "SELECT * FROM vouchers WHERE id = ?", id)
	if err == sql.ErrNoRows { return nil, nil }
	if err != nil { return nil, err }
	return &v, nil
}

func (r *VoucherRepository) FindByCode(ctx context.Context, code string) (*domain.Voucher, error) {
	var v domain.Voucher
	err := r.db.GetContext(ctx, &v, "SELECT * FROM vouchers WHERE code = ?", code)
	if err == sql.ErrNoRows { return nil, nil }
	if err != nil { return nil, err }
	return &v, nil
}

func (r *VoucherRepository) Create(ctx context.Context, v *domain.Voucher) error {
	query := `INSERT INTO vouchers (id, code, title, description, type, discount_value, min_spend, max_discount, usage_limit, per_user_limit, scope, shop_id, category_id, sku, region, payment_method, start_time, end_time, status, stackable, priority, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, v.ID, v.Code, v.Title, v.Description, v.Type, v.DiscountValue, v.MinSpend, v.MaxDiscount, v.UsageLimit, v.PerUserLimit, v.Scope, v.ShopID, v.CategoryID, v.SKU, v.Region, v.PaymentMethod, v.StartTime, v.EndTime, v.Status, v.Stackable, v.Priority, v.CreatedAt, v.UpdatedAt)
	return err
}

func (r *VoucherRepository) Update(ctx context.Context, v *domain.Voucher) error {
	query := `UPDATE vouchers SET title = ?, description = ?, discount_value = ?, min_spend = ?, max_discount = ?, usage_limit = ?, per_user_limit = ?, status = ?, stackable = ?, priority = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, v.Title, v.Description, v.DiscountValue, v.MinSpend, v.MaxDiscount, v.UsageLimit, v.PerUserLimit, v.Status, v.Stackable, v.Priority, v.UpdatedAt, v.ID)
	return err
}

func (r *VoucherRepository) IncrementUsage(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE vouchers SET usage_count = usage_count + 1, updated_at = NOW() WHERE id = ?", id)
	return err
}

func (r *VoucherRepository) ListActive(ctx context.Context, offset, limit int) ([]*domain.Voucher, int64, error) {
	if limit < 1 { limit = 20 }
	if limit > 200 { limit = 200 }
	var total int64
	r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM vouchers WHERE status = 'active' AND end_time > NOW()")
	var vouchers []*domain.Voucher
	err := r.db.SelectContext(ctx, &vouchers, "SELECT * FROM vouchers WHERE status = 'active' AND end_time > NOW() ORDER BY priority DESC LIMIT ? OFFSET ?", limit, offset)
	return vouchers, total, err
}

type VoucherRedemptionRepository struct{ db *sqlx.DB }
func NewVoucherRedemptionRepository(db *sqlx.DB) *VoucherRedemptionRepository {
	return &VoucherRedemptionRepository{db: db}
}

func (r *VoucherRedemptionRepository) FindByID(ctx context.Context, id string) (*domain.VoucherRedemption, error) {
	var red domain.VoucherRedemption
	err := r.db.GetContext(ctx, &red, "SELECT * FROM voucher_redemptions WHERE id = ?", id)
	if err == sql.ErrNoRows { return nil, nil }
	return &red, err
}

func (r *VoucherRedemptionRepository) FindByUserAndVoucher(ctx context.Context, userID, voucherID string) ([]*domain.VoucherRedemption, error) {
	var reds []*domain.VoucherRedemption
	err := r.db.SelectContext(ctx, &reds, "SELECT * FROM voucher_redemptions WHERE user_id = ? AND voucher_id = ?", userID, voucherID)
	return reds, err
}

func (r *VoucherRedemptionRepository) FindByIdempotencyKey(ctx context.Context, key string) (*domain.VoucherRedemption, error) {
	var red domain.VoucherRedemption
	err := r.db.GetContext(ctx, &red, "SELECT * FROM voucher_redemptions WHERE idempotency_key = ?", key)
	if err == sql.ErrNoRows { return nil, nil }
	return &red, err
}

func (r *VoucherRedemptionRepository) CountByVoucher(ctx context.Context, voucherID string) (int64, error) {
	var count int64
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM voucher_redemptions WHERE voucher_id = ?", voucherID)
	return count, err
}

func (r *VoucherRedemptionRepository) CountByUserAndVoucher(ctx context.Context, userID, voucherID string) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM voucher_redemptions WHERE user_id = ? AND voucher_id = ?", userID, voucherID)
	return count, err
}

func (r *VoucherRedemptionRepository) Create(ctx context.Context, red *domain.VoucherRedemption) error {
	query := `INSERT INTO voucher_redemptions (id, voucher_id, user_id, order_id, discount_amount, idempotency_key, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, red.ID, red.VoucherID, red.UserID, red.OrderID, red.DiscountAmt, red.IdempotencyKey, red.CreatedAt)
	return err
}

type CampaignRepository struct{ db *sqlx.DB }
func NewCampaignRepository(db *sqlx.DB) *CampaignRepository { return &CampaignRepository{db: db} }

func (r *CampaignRepository) FindByID(ctx context.Context, id string) (*domain.Campaign, error) {
	var c domain.Campaign
	err := r.db.GetContext(ctx, &c, "SELECT * FROM campaigns WHERE id = ?", id)
	if err == sql.ErrNoRows { return nil, nil }
	return &c, err
}

func (r *CampaignRepository) Create(ctx context.Context, c *domain.Campaign) error {
	query := `INSERT INTO campaigns (id, name, type, description, rules, start_time, end_time, status, priority, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.Name, c.Type, c.Description, c.Rules, c.StartTime, c.EndTime, c.Status, c.Priority, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *CampaignRepository) Update(ctx context.Context, c *domain.Campaign) error {
	query := `UPDATE campaigns SET name = ?, description = ?, rules = ?, status = ?, priority = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, c.Name, c.Description, c.Rules, c.Status, c.Priority, c.UpdatedAt, c.ID)
	return err
}

func (r *CampaignRepository) ListActive(ctx context.Context) ([]*domain.Campaign, error) {
	var campaigns []*domain.Campaign
	err := r.db.SelectContext(ctx, &campaigns, "SELECT * FROM campaigns WHERE status = 'active' AND start_time <= NOW() AND end_time > NOW() ORDER BY priority DESC LIMIT 100")
	return campaigns, err
}

func (r *CampaignRepository) ListByType(ctx context.Context, cType string, offset, limit int) ([]*domain.Campaign, int64, error) {
	if limit < 1 { limit = 20 }
	if limit > 200 { limit = 200 }
	var total int64
	r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM campaigns WHERE type = ?", cType)
	var campaigns []*domain.Campaign
	err := r.db.SelectContext(ctx, &campaigns, "SELECT * FROM campaigns WHERE type = ? ORDER BY priority DESC LIMIT ? OFFSET ?", cType, limit, offset)
	return campaigns, total, err
}

type PricingRuleRepository struct{ db *sqlx.DB }
func NewPricingRuleRepository(db *sqlx.DB) *PricingRuleRepository { return &PricingRuleRepository{db: db} }
func (r *PricingRuleRepository) FindByCampaign(ctx context.Context, campaignID string) ([]*domain.PricingRule, error) {
	var rules []*domain.PricingRule
	err := r.db.SelectContext(ctx, &rules, "SELECT * FROM pricing_rules WHERE campaign_id = ? AND is_active = true ORDER BY priority DESC", campaignID)
	return rules, err
}
func (r *PricingRuleRepository) FindByCampaigns(ctx context.Context, campaignIDs []string) (map[string][]*domain.PricingRule, error) {
	if len(campaignIDs) == 0 {
		return map[string][]*domain.PricingRule{}, nil
	}
	var rules []*domain.PricingRule
	query, args, err := sqlx.In("SELECT * FROM pricing_rules WHERE campaign_id IN (?) AND is_active = true ORDER BY priority DESC", campaignIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &rules, query, args...)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]*domain.PricingRule, len(campaignIDs))
	for _, rule := range rules {
		result[rule.CampaignID] = append(result[rule.CampaignID], rule)
	}
	return result, nil
}
func (r *PricingRuleRepository) Create(ctx context.Context, rule *domain.PricingRule) error {
	query := `INSERT INTO pricing_rules (id, campaign_id, rule_type, condition_json, action_json, priority, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, rule.ID, rule.CampaignID, rule.RuleType, rule.Condition, rule.Action, rule.Priority, rule.IsActive, rule.CreatedAt, rule.UpdatedAt)
	return err
}
func (r *PricingRuleRepository) Update(ctx context.Context, rule *domain.PricingRule) error {
	query := `UPDATE pricing_rules SET condition_json = ?, action_json = ?, priority = ?, is_active = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, rule.Condition, rule.Action, rule.Priority, rule.IsActive, rule.UpdatedAt, rule.ID)
	return err
}

type EligibilityRuleRepository struct{ db *sqlx.DB }
func NewEligibilityRuleRepository(db *sqlx.DB) *EligibilityRuleRepository { return &EligibilityRuleRepository{db: db} }
func (r *EligibilityRuleRepository) FindByPromotion(ctx context.Context, promotionID string) ([]*domain.EligibilityRule, error) {
	var rules []*domain.EligibilityRule
	err := r.db.SelectContext(ctx, &rules, "SELECT * FROM eligibility_rules WHERE promotion_id = ? AND is_active = true ORDER BY id ASC LIMIT 100", promotionID)
	return rules, err
}
func (r *EligibilityRuleRepository) Create(ctx context.Context, rule *domain.EligibilityRule) error {
	query := `INSERT INTO eligibility_rules (id, promotion_id, target_type, target_value, is_active) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, rule.ID, rule.PromotionID, rule.TargetType, rule.TargetValue, rule.IsActive)
	return err
}

type StackingRuleRepository struct{ db *sqlx.DB }
func NewStackingRuleRepository(db *sqlx.DB) *StackingRuleRepository { return &StackingRuleRepository{db: db} }
func (r *StackingRuleRepository) FindByPromotionType(ctx context.Context, pType string) ([]*domain.StackingRule, error) {
	var rules []*domain.StackingRule
	err := r.db.SelectContext(ctx, &rules, "SELECT * FROM stacking_rules WHERE promotion_type = ? ORDER BY id ASC LIMIT 100", pType)
	return rules, err
}
func (r *StackingRuleRepository) Create(ctx context.Context, rule *domain.StackingRule) error {
	query := `INSERT INTO stacking_rules (id, promotion_type, can_stack_with, max_stack_count, priority) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, rule.ID, rule.PromotionType, rule.CanStackWith, rule.MaxStackCount, rule.Priority)
	return err
}
