-- ============================================================
-- Migration: Cart Service Performance Indexes
-- Adds covering indexes for hot-path queries
-- ============================================================

-- Covering index for cart lookup by user (avoids table access)
ALTER TABLE carts ADD INDEX idx_carts_user_cover (user_id, status, deleted_at, id, item_count, subtotal, version, expires_at);

-- Covering index for cart lookup by session
ALTER TABLE carts ADD INDEX idx_carts_session_cover (session_id, status, deleted_at, id, item_count, subtotal, version, expires_at);

-- Covering index for expired cart cleanup
ALTER TABLE carts ADD INDEX idx_carts_expiry_cover (expires_at, status, deleted_at, id);

-- Covering index for cart items listing
ALTER TABLE cart_items ADD INDEX idx_cart_items_cover (cart_id, added_at DESC, sku, product_name, quantity, unit_price, total_price, is_selected, is_available);

-- Covering index for cart snapshot lookup
ALTER TABLE cart_snapshots ADD INDEX idx_snapshots_cart_cover (cart_id, expires_at, idempotency_key);

ANALYZE TABLE carts;
ANALYZE TABLE cart_items;
ANALYZE TABLE cart_snapshots;
