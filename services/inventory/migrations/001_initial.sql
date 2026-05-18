CREATE TABLE IF NOT EXISTS stock (
    id VARCHAR(36) PRIMARY KEY, product_id VARCHAR(36) NOT NULL, sku_id VARCHAR(36) NOT NULL,
    warehouse_id VARCHAR(36) NOT NULL, quantity INT NOT NULL DEFAULT 0,
    reserved_qty INT NOT NULL DEFAULT 0, available_qty INT NOT NULL DEFAULT 0,
    status ENUM('in_stock','low_stock','out_of_stock','reserved') NOT NULL DEFAULT 'in_stock',
    reorder_level INT NOT NULL DEFAULT 10, version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_sku_warehouse (sku_id, warehouse_id),
    INDEX idx_stock_product (product_id), INDEX idx_stock_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS reservations (
    id VARCHAR(36) PRIMARY KEY, order_id VARCHAR(36) NOT NULL, user_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL, sku_id VARCHAR(36) NOT NULL, warehouse_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL DEFAULT 0, status ENUM('active','committed','released','expired') NOT NULL DEFAULT 'active',
    expires_at TIMESTAMP NOT NULL, idempotency_key VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_reservations_order (order_id), INDEX idx_reservations_sku (sku_id),
    INDEX idx_reservations_status (status), INDEX idx_reservations_expires (expires_at),
    INDEX idx_reservations_idempotency (idempotency_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS idempotency_keys (
    ` + "`key`" + ` VARCHAR(255) PRIMARY KEY, reservation_id VARCHAR(36) NOT NULL,
    expires_at TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_inv_idempotency_reservation (reservation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS outbox_events (
    event_id VARCHAR(36) PRIMARY KEY, aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL, event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    INDEX idx_inv_outbox_processed (processed, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
