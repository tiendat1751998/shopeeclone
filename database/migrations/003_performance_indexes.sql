-- ============================================================
-- Migration 003: Performance Indexes
-- Adds missing composite indexes identified by query analysis
-- ============================================================

-- 1. Products: composite index for search + sort by creation
--    WHERE status='active' AND deleted_at IS NULL ORDER BY created_at DESC
ALTER TABLE products ADD INDEX idx_prod_status_created (status, deleted_at, created_at);

-- 2. SKUs: index for price filtering and sorting
ALTER TABLE skus ADD INDEX idx_sku_sale_price (sale_price);
ALTER TABLE skus ADD INDEX idx_sku_product_price (product_id, price);

-- 3. Product media: composite for sorting
ALTER TABLE product_media ADD INDEX idx_pm_product_sort (product_id, sort_order);

-- 4. Sessions: composite for active session lookups
ALTER TABLE sessions ADD INDEX idx_sessions_user_status (user_id, status, expires_at);

-- 5. Categories: composite for slug lookup
ALTER TABLE categories ADD INDEX idx_cat_slug_active (slug, is_active);
