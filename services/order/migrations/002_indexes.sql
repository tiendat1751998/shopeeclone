-- Order Service - Additional Indexes & Optimizations
-- Migration: 002_indexes

-- Composite index for common query patterns
CREATE INDEX idx_orders_user_status_date ON orders(user_id, status, created_at DESC);
CREATE INDEX idx_orders_seller_status_date ON orders(seller_id, status, created_at DESC);

-- Covering index for order listing
CREATE INDEX idx_orders_list_covering ON orders(user_id, status, created_at DESC, id, order_number, total_amount, currency);

-- Cleanup expired idempotency keys (run periodically)
-- DELETE FROM idempotency_keys WHERE expires_at < NOW();

-- Cleanup soft-deleted orders older than retention period (run periodically)
-- DELETE FROM orders WHERE deleted_at IS NOT NULL AND deleted_at < DATE_SUB(NOW(), INTERVAL 90 DAY);
