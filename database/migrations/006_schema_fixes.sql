-- ============================================================
-- Migration 006: Fix schema mismatches for categories, skus, products
-- Adds missing columns that the Go and Next.js code expect
-- ============================================================

-- Fix categories table: add missing columns
ALTER TABLE categories
  ADD COLUMN IF NOT EXISTS description TEXT DEFAULT NULL AFTER slug,
  ADD COLUMN IF NOT EXISTS image_url VARCHAR(500) DEFAULT NULL AFTER description,
  ADD COLUMN IF NOT EXISTS depth INT NOT NULL DEFAULT 0 AFTER is_active,
  ADD COLUMN IF NOT EXISTS path VARCHAR(1000) DEFAULT NULL AFTER depth,
  ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 1 AFTER path;

-- Fix skus table: add missing columns
ALTER TABLE skus
  ADD COLUMN IF NOT EXISTS name VARCHAR(500) DEFAULT NULL AFTER sku_code,
  ADD COLUMN IF NOT EXISTS compare_price BIGINT DEFAULT NULL AFTER price,
  ADD COLUMN IF NOT EXISTS currency VARCHAR(3) NOT NULL DEFAULT 'VND' AFTER compare_price,
  ADD COLUMN IF NOT EXISTS reserved_stock BIGINT NOT NULL DEFAULT 0 AFTER stock,
  ADD COLUMN IF NOT EXISTS weight_grams INT DEFAULT NULL AFTER reserved_stock,
  ADD COLUMN IF NOT EXISTS dimensions VARCHAR(255) DEFAULT NULL AFTER weight_grams,
  ADD COLUMN IF NOT EXISTS metadata TEXT DEFAULT NULL AFTER attributes,
  ADD COLUMN IF NOT EXISTS sort_order INT NOT NULL DEFAULT 0 AFTER metadata,
  ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 1 AFTER sort_order;

-- Fix products table: add missing columns
ALTER TABLE products
  ADD COLUMN IF NOT EXISTS brand VARCHAR(255) DEFAULT NULL AFTER category_id,
  ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(255) DEFAULT NULL AFTER currency,
  ADD INDEX IF NOT EXISTS idx_products_idempotency (idempotency_key);
