-- ============================================================
-- Migration: Promotion Service Performance Indexes
-- Adds covering indexes for hot-path queries
-- ============================================================

-- Covering index for voucher validation (the hottest promotion query)
ALTER TABLE vouchers ADD INDEX idx_vouchers_code_cover (code, status, start_time, end_time, usage_limit, usage_count, per_user_limit, scope, shop_id, category_id, sku, discount_value, min_spend, max_discount);

-- Covering index for active voucher listing
ALTER TABLE vouchers ADD INDEX idx_vouchers_active_cover (status, start_time, end_time, scope, shop_id, id, code, discount_value, min_spend);

-- Covering index for voucher redemptions by user+voucher
ALTER TABLE voucher_redemptions ADD INDEX idx_redemptions_user_voucher_cover (user_id, voucher_id, created_at DESC, id, order_id, discount_amount);

-- Covering index for active campaigns lookup
ALTER TABLE campaigns ADD INDEX idx_campaigns_active_cover (status, start_time, end_time, type, priority DESC, id, name);

ANALYZE TABLE vouchers;
ANALYZE TABLE voucher_redemptions;
ANALYZE TABLE campaigns;
