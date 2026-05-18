-- Checkout Service Migration 001

CREATE TABLE IF NOT EXISTS checkouts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    cart_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) DEFAULT NULL,
    status ENUM('pending','validating','pricing_frozen','reserving_inventory','inventory_reserved','processing_payment','completed','failed','rolling_back','rolled_back','expired') NOT NULL DEFAULT 'pending',
    idempotency_key VARCHAR(100) DEFAULT NULL,
    current_step VARCHAR(50) NOT NULL DEFAULT 'init',
    failure_reason TEXT DEFAULT NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    reservation_keys TEXT DEFAULT NULL,
    pricing_snapshot VARCHAR(36) DEFAULT NULL,
    promotion_results TEXT DEFAULT NULL,
    subtotal BIGINT NOT NULL DEFAULT 0,
    discount_total BIGINT NOT NULL DEFAULT 0,
    shipping_total BIGINT NOT NULL DEFAULT 0,
    grand_total BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    expires_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_checkouts_user (user_id, status),
    INDEX idx_checkouts_idempotency (idempotency_key),
    INDEX idx_checkouts_status (status, expires_at),
    INDEX idx_checkouts_cart (cart_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS checkout_step_logs (
    id VARCHAR(36) PRIMARY KEY,
    checkout_id VARCHAR(36) NOT NULL,
    step VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_message TEXT DEFAULT NULL,
    metadata TEXT DEFAULT NULL,
    duration_ms BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_step_logs_checkout (checkout_id),
    FOREIGN KEY (checkout_id) REFERENCES checkouts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS pricing_snapshots (
    id VARCHAR(36) PRIMARY KEY,
    checkout_id VARCHAR(36) NOT NULL,
    items JSON NOT NULL,
    seller_groups JSON NOT NULL,
    subtotal BIGINT NOT NULL,
    discount_total BIGINT NOT NULL DEFAULT 0,
    shipping_total BIGINT NOT NULL DEFAULT 0,
    grand_total BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    promotions_applied JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_pricing_checkout (checkout_id),
    FOREIGN KEY (checkout_id) REFERENCES checkouts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS reservation_orchestrations (
    id VARCHAR(36) PRIMARY KEY,
    checkout_id VARCHAR(36) NOT NULL,
    reservation_key VARCHAR(100) NOT NULL,
    sku VARCHAR(100) NOT NULL,
    warehouse_id VARCHAR(36) NOT NULL,
    quantity BIGINT NOT NULL,
    status ENUM('pending','reserved','released','failed') NOT NULL DEFAULT 'pending',
    error_message TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_res_orchestration_checkout (checkout_id),
    INDEX idx_res_orchestration_key (reservation_key),
    FOREIGN KEY (checkout_id) REFERENCES checkouts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS reconciliation_jobs (
    id VARCHAR(36) PRIMARY KEY,
    checkout_id VARCHAR(36) NOT NULL,
    job_type ENUM('release_reservation','confirm_reservation','update_order_status') NOT NULL,
    status ENUM('pending','running','completed','failed') NOT NULL DEFAULT 'pending',
    attempt_count INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    next_retry_at TIMESTAMP NOT NULL,
    metadata TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_reconcile_status (status, next_retry_at),
    FOREIGN KEY (checkout_id) REFERENCES checkouts(id)
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
