CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY, shop_id VARCHAR(36) NOT NULL, name VARCHAR(500) NOT NULL, description TEXT,
    category_id VARCHAR(36) NOT NULL, status ENUM('draft','active','inactive','archived','moderated') NOT NULL DEFAULT 'draft',
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD', version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_products_shop (shop_id, status, deleted_at), INDEX idx_products_category (category_id, status, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS skus (
    id VARCHAR(36) PRIMARY KEY, product_id VARCHAR(36) NOT NULL, sku_code VARCHAR(100) NOT NULL, attributes TEXT,
    price BIGINT NOT NULL, sale_price BIGINT DEFAULT NULL, stock BIGINT NOT NULL DEFAULT 0,
    status ENUM('active','inactive','out_of_stock') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_skus_product (product_id), UNIQUE KEY uk_product_sku (product_id, sku_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(36) PRIMARY KEY, parent_id VARCHAR(36) DEFAULT NULL, name VARCHAR(255) NOT NULL, slug VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0, sort_order INT NOT NULL DEFAULT 0, is_active BOOLEAN DEFAULT TRUE, metadata TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_categories_parent (parent_id, is_active), INDEX idx_categories_level (level, sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS attributes (
    id VARCHAR(36) PRIMARY KEY, category_id VARCHAR(36) NOT NULL, name VARCHAR(100) NOT NULL, display_name VARCHAR(255) NOT NULL,
    type ENUM('text','number','select','multi_select','boolean') NOT NULL, required BOOLEAN DEFAULT FALSE, options TEXT, sort_order INT DEFAULT 0, is_active BOOLEAN DEFAULT TRUE,
    INDEX idx_attrs_category (category_id, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS product_media (
    id VARCHAR(36) PRIMARY KEY, product_id VARCHAR(36) NOT NULL, media_type ENUM('image','video') NOT NULL,
    url VARCHAR(500) NOT NULL, thumbnail VARCHAR(500) DEFAULT NULL, sort_order INT DEFAULT 0, is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, INDEX idx_media_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS outbox_events (
    event_id VARCHAR(36) PRIMARY KEY, aggregate_type VARCHAR(100) NOT NULL, aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL, payload JSON NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, processed BOOLEAN DEFAULT FALSE,
    INDEX idx_outbox_processed (processed, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
