-- ============================================================
-- TIKI FEATURES DATABASE SCHEMA
-- Adds TikiNow, TikiXu, Tiki Trading features
-- ============================================================

-- TikiNow - Express Delivery Zones
CREATE TABLE IF NOT EXISTS tiki_now_zones (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    ward VARCHAR(100) DEFAULT NULL,
    max_delivery_minutes INT NOT NULL DEFAULT 120,
    available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tnz_city (city, available),
    INDEX idx_tnz_district (district)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Products eligible for TikiNow
CREATE TABLE IF NOT EXISTS tiki_now_products (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL UNIQUE,
    sku_id VARCHAR(36) NOT NULL,
    warehouse_id VARCHAR(36) NOT NULL,
    max_quantity INT NOT NULL DEFAULT 10,
    cutoff_time TIME NOT NULL DEFAULT '16:00:00',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tnp_product (product_id),
    INDEX idx_tnp_active (is_active),
    INDEX idx_tnp_warehouse (warehouse_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- TikiXu - Loyalty Points
CREATE TABLE IF NOT EXISTS tiki_xu_accounts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    balance BIGINT NOT NULL DEFAULT 0,
    lifetime_earned BIGINT NOT NULL DEFAULT 0,
    lifetime_spent BIGINT NOT NULL DEFAULT 0,
    tier ENUM('bronze','silver','gold','platinum','diamond') NOT NULL DEFAULT 'bronze',
    tier_points INT NOT NULL DEFAULT 0,
    status ENUM('active','frozen','closed') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_txa_tier (tier),
    INDEX idx_txa_status (status),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tiki_xu_transactions (
    id VARCHAR(36) PRIMARY KEY,
    account_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    type ENUM('earn','spend','refund','expire','adjust','signup_bonus','purchase','review','referral') NOT NULL,
    reference_type VARCHAR(50) DEFAULT NULL,
    reference_id VARCHAR(36) DEFAULT NULL,
    description VARCHAR(500) DEFAULT NULL,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_txt_account (account_id),
    INDEX idx_txt_user (user_id, created_at),
    INDEX idx_txt_type (type),
    INDEX idx_txt_reference (reference_type, reference_id),
    FOREIGN KEY (account_id) REFERENCES tiki_xu_accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tiki_xu_tier_benefits (
    id VARCHAR(36) PRIMARY KEY,
    tier ENUM('bronze','silver','gold','platinum','diamond') NOT NULL UNIQUE,
    min_points INT NOT NULL DEFAULT 0,
    earn_rate DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    free_shipping BOOLEAN NOT NULL DEFAULT FALSE,
    exclusive_deals BOOLEAN NOT NULL DEFAULT FALSE,
    priority_support BOOLEAN NOT NULL DEFAULT FALSE,
    birthday_gift BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tiki Trading (first-party products)
ALTER TABLE products ADD COLUMN IF NOT EXISTS is_tiki_trading BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE products ADD COLUMN IF NOT EXISTS tiki_trading_type ENUM('direct_import','marketplace','cross_border') DEFAULT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS tiki_guarantee_months INT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN IF NOT EXISTS tiki_original_price BIGINT DEFAULT NULL;

-- TikiNow orders extension
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS is_tiki_now BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS delivery_promise_minutes INT DEFAULT NULL;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS delivery_deadline TIMESTAMP NULL DEFAULT NULL;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS tiki_now_fee BIGINT NOT NULL DEFAULT 0;

-- TikiXu applied to checkout
ALTER TABLE checkouts ADD COLUMN IF NOT EXISTS tiki_xu_applied BIGINT NOT NULL DEFAULT 0;
ALTER TABLE checkouts ADD COLUMN IF NOT EXISTS tiki_xu_discount BIGINT NOT NULL DEFAULT 0;
ALTER TABLE checkouts ADD COLUMN IF NOT EXISTS tiki_xu_rate DECIMAL(5,2) DEFAULT NULL;

-- Product TikiNow eligibility
ALTER TABLE products ADD COLUMN IF NOT EXISTS is_tiki_now_eligible BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE products ADD COLUMN IF NOT EXISTS tiki_now_cutoff_time TIME DEFAULT NULL;
ALTER TABLE products ADD COLUMN IF NOT EXISTS tiki_now_max_qty INT NOT NULL DEFAULT 0;

-- Review media (Tiki-style reviews with images)
CREATE TABLE IF NOT EXISTS review_media (
    id VARCHAR(36) PRIMARY KEY,
    review_id VARCHAR(36) NOT NULL,
    url VARCHAR(500) NOT NULL,
    media_type ENUM('image','video') NOT NULL DEFAULT 'image',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_rm_review (review_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tiki Flash Sale / Mega Deals
CREATE TABLE IF NOT EXISTS tiki_deals (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    banner_url VARCHAR(500) DEFAULT NULL,
    deal_type ENUM('flash_sale','mega_deal','daily_discover','brand_sale') NOT NULL DEFAULT 'flash_sale',
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    status ENUM('draft','active','paused','ended') NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_td_type (deal_type, status),
    INDEX idx_td_time (start_time, end_time, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tiki_deal_products (
    id VARCHAR(36) PRIMARY KEY,
    deal_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    sku_id VARCHAR(36) NOT NULL,
    deal_price BIGINT NOT NULL,
    original_price BIGINT NOT NULL,
    max_quantity INT NOT NULL DEFAULT 100,
    sold_quantity INT NOT NULL DEFAULT 0,
    max_per_user INT NOT NULL DEFAULT 1,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tdp_deal (deal_id),
    INDEX idx_tdp_product (product_id),
    UNIQUE KEY uk_tdp_deal_sku (deal_id, sku_id),
    FOREIGN KEY (deal_id) REFERENCES tiki_deals(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tiki Xu tier benefits seed
INSERT IGNORE INTO tiki_xu_tier_benefits (id, tier, min_points, earn_rate, free_shipping, exclusive_deals, priority_support, birthday_gift, description) VALUES
('txb-001', 'bronze', 0, 1.00, FALSE, FALSE, FALSE, FALSE, 'Basic membership - earn 1% back on every purchase'),
('txb-002', 'silver', 1000, 1.25, FALSE, FALSE, FALSE, TRUE, 'Silver - earn 1.25% back, birthday gift'),
('txb-003', 'gold', 5000, 1.50, TRUE, FALSE, FALSE, TRUE, 'Gold - earn 1.5% back, free shipping, birthday gift'),
('txb-004', 'platinum', 20000, 2.00, TRUE, TRUE, TRUE, TRUE, 'Platinum - earn 2% back, free shipping, exclusive deals'),
('txb-005', 'diamond', 50000, 3.00, TRUE, TRUE, TRUE, TRUE, 'Diamond - earn 3% back, all perks unlocked');

-- Tiki product reviews with media
CREATE TABLE IF NOT EXISTS tiki_reviews (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    sku_id VARCHAR(36) DEFAULT NULL,
    order_id VARCHAR(36) DEFAULT NULL,
    rating TINYINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255) DEFAULT NULL,
    content TEXT DEFAULT NULL,
    is_verified_purchase BOOLEAN NOT NULL DEFAULT FALSE,
    is_recommended BOOLEAN NOT NULL DEFAULT TRUE,
    likes_count INT NOT NULL DEFAULT 0,
    status ENUM('pending','approved','rejected','flagged') NOT NULL DEFAULT 'approved',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tr_product (product_id, status),
    INDEX idx_tr_user (user_id),
    INDEX idx_tr_rating (rating),
    INDEX idx_tr_created (created_at),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
