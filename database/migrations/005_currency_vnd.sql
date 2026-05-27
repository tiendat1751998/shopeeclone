-- ============================================================
-- Migration 005: Currency Configuration Update
-- Change platform currency from SGD to VND
-- VND is the primary currency, no decimal places needed
-- BIGINT prices are in whole VND (no cents/satang)
-- ============================================================

-- Update default currency in all tables
ALTER TABLE products   MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE orders     MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE order_items MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE carts      MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE cart_items MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE cart_snapshots MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE payments   MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE refunds    MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE checkouts  MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';
ALTER TABLE pricing_snapshots MODIFY COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'VND';

-- Update existing data to VND (only if currently SGD and no other currency set)
-- WARNING: Only run this if your data is still in SGD and you want to convert
-- If you want to keep existing data as-is, skip these UPDATE statements
UPDATE products SET currency = 'VND' WHERE currency = 'SGD';
UPDATE orders SET currency = 'VND' WHERE currency = 'SGD';
UPDATE payments SET currency = 'VND' WHERE currency = 'SGD';
