-- Cart Service Migration 001

CREATE TABLE IF NOT EXISTS carts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) DEFAULT NULL,
    session_id VARCHAR(100) DEFAULT NULL,
    status ENUM('active', 'merged', 'expired', 'checkout') NOT NULL DEFAULT 'active',
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    item_count INT NOT NULL DEFAULT 0,
    subtotal BIGINT NOT NULL DEFAULT 0,
    version BIGINT NOT NULL DEFAULT 1,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_carts_user (user_id, status, deleted_at),
    INDEX idx_carts_session (session_id, status, deleted_at),
    INDEX idx_carts_expires (expires_at, status),
    UNIQUE INDEX idx_uq_active_user_cart (user_id) WHERE status = 'active' AND deleted_at IS NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS cart_items (
    id VARCHAR(36) PRIMARY KEY,
    cart_id VARCHAR(36) NOT NULL,
    sku VARCHAR(100) NOT NULL,
    product_name VARCHAR(500) NOT NULL,
    shop_id VARCHAR(36) NOT NULL,
    shop_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price BIGINT NOT NULL,
    total_price BIGINT NOT NULL,
    image_url VARCHAR(500) DEFAULT NULL,
    attributes TEXT DEFAULT NULL,
    is_selected BOOLEAN DEFAULT TRUE,
    is_available BOOLEAN DEFAULT TRUE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_cart_items_cart (cart_id),
    INDEX idx_cart_items_sku (sku),
    INDEX idx_cart_items_shop (shop_id),
    UNIQUE KEY uk_cart_sku (cart_id, sku),
    FOREIGN KEY (cart_id) REFERENCES carts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS cart_snapshots (
    id VARCHAR(36) PRIMARY KEY,
    cart_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    items JSON NOT NULL,
    seller_groups JSON NOT NULL,
    subtotal BIGINT NOT NULL,
    item_count INT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    idempotency_key VARCHAR(100) DEFAULT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_snapshots_cart (cart_id),
    INDEX idx_snapshots_idempotency (idempotency_key),
    INDEX idx_snapshots_expires (expires_at),
    FOREIGN KEY (cart_id) REFERENCES carts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS cart_merge_history (
    id VARCHAR(36) PRIMARY KEY,
    source_cart_id VARCHAR(36) NOT NULL,
    target_cart_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    merge_type ENUM('guest_to_user', 'session', 'conflict_resolution') NOT NULL,
    items_merged INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_merge_user (user_id),
    INDEX idx_merge_target (target_cart_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS outbox_events (
    event_id VARCHAR(36) PRIMARY KEY,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE,
    INDEX idx_outbox_processed (processed, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
