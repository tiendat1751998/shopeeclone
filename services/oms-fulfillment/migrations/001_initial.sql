-- OMS Fulfillment Service - Initial Schema
-- Order management, fulfillment workflows, and warehouse operations

CREATE TABLE IF NOT EXISTS fulfillment_orders (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    seller_id VARCHAR(36) NOT NULL,
    warehouse_id VARCHAR(36) DEFAULT NULL,
    status ENUM('pending','confirmed','picking','packed','shipped','delivered','returned','cancelled','failed') NOT NULL DEFAULT 'pending',
    priority ENUM('low','normal','high','urgent') NOT NULL DEFAULT 'normal',
    shipping_method VARCHAR(64) NOT NULL,
    shipping_address JSON NOT NULL,
    estimated_ship_date TIMESTAMP NULL DEFAULT NULL,
    estimated_delivery_date TIMESTAMP NULL DEFAULT NULL,
    actual_shipped_at TIMESTAMP NULL DEFAULT NULL,
    actual_delivered_at TIMESTAMP NULL DEFAULT NULL,
    tracking_number VARCHAR(255) DEFAULT NULL,
    carrier VARCHAR(64) DEFAULT NULL,
    total_items INT NOT NULL DEFAULT 0,
    total_weight_grams INT NOT NULL DEFAULT 0,
    notes TEXT DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_fo_order (order_id),
    INDEX idx_fo_seller (seller_id, status),
    INDEX idx_fo_warehouse (warehouse_id, status),
    INDEX idx_fo_status (status),
    INDEX idx_fo_tracking (tracking_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS fulfillment_items (
    id VARCHAR(36) PRIMARY KEY,
    fulfillment_id VARCHAR(36) NOT NULL,
    order_item_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    sku_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    picked_quantity INT NOT NULL DEFAULT 0,
    packed_quantity INT NOT NULL DEFAULT 0,
    status ENUM('pending','picked','packed','shipped','returned','cancelled') NOT NULL DEFAULT 'pending',
    picked_at TIMESTAMP NULL DEFAULT NULL,
    packed_at TIMESTAMP NULL DEFAULT NULL,
    picked_by VARCHAR(36) DEFAULT NULL,
    packed_by VARCHAR(36) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_fi_fulfillment (fulfillment_id),
    INDEX idx_fi_product (product_id),
    INDEX idx_fi_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS warehouses (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(36) NOT NULL UNIQUE,
    address JSON NOT NULL,
    contact_phone VARCHAR(20) DEFAULT NULL,
    contact_email VARCHAR(255) DEFAULT NULL,
    capacity_total INT NOT NULL DEFAULT 0,
    capacity_used INT NOT NULL DEFAULT 0,
    operating_hours JSON DEFAULT NULL,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_wh_code (code),
    INDEX idx_wh_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS warehouse_inventory (
    id VARCHAR(36) PRIMARY KEY,
    warehouse_id VARCHAR(36) NOT NULL,
    sku_id VARCHAR(36) NOT NULL,
    quantity_available INT NOT NULL DEFAULT 0,
    quantity_reserved INT NOT NULL DEFAULT 0,
    quantity_damaged INT NOT NULL DEFAULT 0,
    reorder_level INT NOT NULL DEFAULT 0,
    reorder_quantity INT NOT NULL DEFAULT 0,
    location_code VARCHAR(64) DEFAULT NULL,
    last_counted_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_whi_wh_sku (warehouse_id, sku_id),
    INDEX idx_whi_sku (sku_id),
    INDEX idx_whi_reorder (quantity_available, reorder_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS fulfillment_events (
    id VARCHAR(36) PRIMARY KEY,
    fulfillment_id VARCHAR(36) NOT NULL,
    event_type ENUM('created','confirmed','pick_started','pick_completed','pack_started','pack_completed','shipped','in_transit','delivered','returned','cancelled','exception') NOT NULL,
    description TEXT DEFAULT NULL,
    actor_id VARCHAR(36) DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_fe_fulfillment (fulfillment_id, created_at),
    INDEX idx_fe_type (event_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS return_requests (
    id VARCHAR(36) PRIMARY KEY,
    fulfillment_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    return_type ENUM('full','partial','exchange') NOT NULL DEFAULT 'partial',
    reason ENUM('defective','wrong_item','not_as_described','changed_mind','damaged_in_shipping','late_delivery','other') NOT NULL,
    reason_detail TEXT DEFAULT NULL,
    status ENUM('requested','approved','rejected','items_shipped','items_received','refunded','completed','cancelled') NOT NULL DEFAULT 'requested',
    refund_amount BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    tracking_number VARCHAR(255) DEFAULT NULL,
    approved_by VARCHAR(36) DEFAULT NULL,
    approved_at TIMESTAMP NULL DEFAULT NULL,
    received_at TIMESTAMP NULL DEFAULT NULL,
    refunded_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_rr_fulfillment (fulfillment_id),
    INDEX idx_rr_order (order_id),
    INDEX idx_rr_user (user_id),
    INDEX idx_rr_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS return_items (
    id VARCHAR(36) PRIMARY KEY,
    return_id VARCHAR(36) NOT NULL,
    fulfillment_item_id VARCHAR(36) NOT NULL,
    sku_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    `condition` ENUM('unopened','opened_unused','used','damaged','defective') NOT NULL DEFAULT 'unopened',
    refund_amount BIGINT NOT NULL DEFAULT 0,
    restockable BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ri_return (return_id),
    INDEX idx_ri_sku (sku_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
