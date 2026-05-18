CREATE TABLE IF NOT EXISTS shipments (
    id VARCHAR(36) PRIMARY KEY, order_id VARCHAR(36) NOT NULL, user_id VARCHAR(36) NOT NULL,
    carrier_id VARCHAR(64) NOT NULL, tracking_number VARCHAR(255) DEFAULT '',
    status ENUM('pending','booked','picked_up','in_transit','out_for_delivery','delivered','failed','returned','cancelled') NOT NULL DEFAULT 'pending',
    origin_address JSON, destination_address JSON, weight DOUBLE NOT NULL DEFAULT 0,
    dimensions VARCHAR(64) DEFAULT '', label_url VARCHAR(512) DEFAULT '',
    cost BIGINT NOT NULL DEFAULT 0, currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    idempotency_key VARCHAR(255) DEFAULT '', metadata JSON, version INT NOT NULL DEFAULT 1,
    estimated_delivery TIMESTAMP NULL, delivered_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_shipments_order_id (order_id), INDEX idx_shipments_user_id (user_id),
    INDEX idx_shipments_status (status), INDEX idx_shipments_tracking (tracking_number),
    INDEX idx_shipments_idempotency (idempotency_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tracking_events (
    id VARCHAR(36) PRIMARY KEY, shipment_id VARCHAR(36) NOT NULL,
    status VARCHAR(64) NOT NULL, location VARCHAR(255) DEFAULT '',
    description TEXT, timestamp TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tracking_shipment_id (shipment_id), INDEX idx_tracking_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS webhook_events (
    id VARCHAR(36) PRIMARY KEY, provider VARCHAR(64) NOT NULL, event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL, signature VARCHAR(255) NOT NULL, idempotency_key VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_shipment_webhooks_idempotency (idempotency_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS idempotency_keys (
    ` + "`key`" + ` VARCHAR(255) PRIMARY KEY, shipment_id VARCHAR(36) NOT NULL,
    expires_at TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_shipment_idempotency_shipment (shipment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS outbox_events (
    event_id VARCHAR(36) PRIMARY KEY, aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL, event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    INDEX idx_shipment_outbox_processed (processed, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
