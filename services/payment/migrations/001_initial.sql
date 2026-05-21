-- Payment Service - Initial Schema

CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    status ENUM('pending','authorized','captured','failed','expired','refunded','partial_refund','cancelled') NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(32) NOT NULL,
    psp_transaction_id VARCHAR(255) DEFAULT '',
    psp_provider VARCHAR(64) NOT NULL,
    idempotency_key VARCHAR(255) DEFAULT '',
    amount_refunded BIGINT NOT NULL DEFAULT 0,
    failure_reason VARCHAR(255) DEFAULT '',
    metadata JSON,
    version INT NOT NULL DEFAULT 1,
    authorized_at TIMESTAMP NULL,
    captured_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_payments_order_id (order_id),
    INDEX idx_payments_user_id (user_id),
    INDEX idx_payments_status (status),
    INDEX idx_payments_idempotency (idempotency_key),
    INDEX idx_payments_psp_txn (psp_transaction_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS refunds (
    id VARCHAR(36) PRIMARY KEY,
    payment_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL,
    status ENUM('pending','processed','failed') NOT NULL DEFAULT 'pending',
    reason TEXT NOT NULL,
    psp_refund_id VARCHAR(255) DEFAULT '',
    idempotency_key VARCHAR(255) DEFAULT '',
    metadata JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_refunds_payment_id (payment_id),
    INDEX idx_refunds_order_id (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS webhook_events (
    id VARCHAR(36) PRIMARY KEY,
    psp_provider VARCHAR(64) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL,
    signature VARCHAR(255) NOT NULL,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    retry_count INT NOT NULL DEFAULT 0,
    idempotency_key VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_webhooks_idempotency (idempotency_key),
    INDEX idx_webhooks_processed (processed)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS fraud_checks (
    id VARCHAR(36) PRIMARY KEY,
    payment_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    risk_score INT NOT NULL DEFAULT 0,
    risk_level VARCHAR(16) NOT NULL DEFAULT 'low',
    is_fraud BOOLEAN NOT NULL DEFAULT FALSE,
    reasons JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_fraud_payment_id (payment_id),
    INDEX idx_fraud_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS idempotency_keys (
    `key` VARCHAR(255) PRIMARY KEY,
    payment_id VARCHAR(36) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_idempotency_payment_id (payment_id),
    INDEX idx_idempotency_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS outbox_events (
    event_id VARCHAR(36) PRIMARY KEY,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    INDEX idx_outbox_processed (processed, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
