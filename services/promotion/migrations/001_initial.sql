-- Promotion Service Migration 001

CREATE TABLE IF NOT EXISTS vouchers (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type ENUM('percentage', 'fixed', 'shipping') NOT NULL,
    discount_value BIGINT NOT NULL,
    min_spend BIGINT NOT NULL DEFAULT 0,
    max_discount BIGINT NOT NULL DEFAULT 0,
    usage_limit BIGINT NOT NULL DEFAULT 10000,
    usage_count BIGINT NOT NULL DEFAULT 0,
    per_user_limit INT NOT NULL DEFAULT 1,
    scope ENUM('platform', 'shop', 'category', 'sku') NOT NULL DEFAULT 'platform',
    shop_id VARCHAR(36) DEFAULT NULL,
    category_id VARCHAR(36) DEFAULT NULL,
    sku VARCHAR(100) DEFAULT NULL,
    region VARCHAR(50) DEFAULT NULL,
    payment_method VARCHAR(50) DEFAULT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status ENUM('active', 'inactive', 'expired', 'exhausted') NOT NULL DEFAULT 'active',
    stackable BOOLEAN DEFAULT FALSE,
    priority INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_vouchers_code (code),
    INDEX idx_vouchers_status (status, start_time, end_time),
    INDEX idx_vouchers_scope (scope, shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS voucher_redemptions (
    id VARCHAR(36) PRIMARY KEY,
    voucher_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) NOT NULL,
    discount_amount BIGINT NOT NULL,
    idempotency_key VARCHAR(100) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_redemptions_voucher (voucher_id),
    INDEX idx_redemptions_user (user_id, voucher_id),
    INDEX idx_redemptions_order (order_id),
    INDEX idx_redemptions_idempotency (idempotency_key),
    FOREIGN KEY (voucher_id) REFERENCES vouchers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS campaigns (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type ENUM('flash_sale', 'scheduled', 'seasonal', 'seller', 'category') NOT NULL,
    description TEXT,
    rules JSON,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status ENUM('draft', 'active', 'paused', 'ended') NOT NULL DEFAULT 'draft',
    priority INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_campaigns_type (type, status),
    INDEX idx_campaigns_time (start_time, end_time, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS pricing_rules (
    id VARCHAR(36) PRIMARY KEY,
    campaign_id VARCHAR(36) NOT NULL,
    rule_type VARCHAR(50) NOT NULL,
    condition_json JSON NOT NULL,
    action_json JSON NOT NULL,
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_pricing_campaign (campaign_id),
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS eligibility_rules (
    id VARCHAR(36) PRIMARY KEY,
    promotion_id VARCHAR(36) NOT NULL,
    target_type ENUM('user', 'region', 'payment', 'seller', 'product') NOT NULL,
    target_value VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    INDEX idx_eligibility_promo (promotion_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS stacking_rules (
    id VARCHAR(36) PRIMARY KEY,
    promotion_type VARCHAR(50) NOT NULL,
    can_stack_with VARCHAR(50) NOT NULL,
    max_stack_count INT DEFAULT 1,
    priority INT DEFAULT 0,
    INDEX idx_stacking_type (promotion_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS outbox_events (
    event_id VARCHAR(36) PRIMARY KEY,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE,
    INDEX idx_outbox_processed (processed, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;
