-- Developer Platform Service - Initial Schema
-- API keys, developer portal, documentation, and usage analytics

CREATE TABLE IF NOT EXISTS developer_accounts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    company_name VARCHAR(255) DEFAULT NULL,
    website_url VARCHAR(500) DEFAULT NULL,
    contact_email VARCHAR(255) NOT NULL,
    status ENUM('pending','active','suspended','banned') NOT NULL DEFAULT 'pending',
    tier ENUM('free','basic','pro','enterprise') NOT NULL DEFAULT 'free',
    verified_at TIMESTAMP NULL DEFAULT NULL,
    approved_by VARCHAR(36) DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_da_user (user_id),
    INDEX idx_da_status (status),
    INDEX idx_da_tier (tier)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_keys (
    id VARCHAR(36) PRIMARY KEY,
    developer_id VARCHAR(36) NOT NULL,
    app_name VARCHAR(255) NOT NULL,
    api_key_hash VARCHAR(255) NOT NULL UNIQUE,
    api_key_prefix VARCHAR(8) NOT NULL,
    scopes JSON NOT NULL,
    rate_limit_per_minute INT NOT NULL DEFAULT 60,
    rate_limit_per_day INT NOT NULL DEFAULT 10000,
    environment ENUM('sandbox','production') NOT NULL DEFAULT 'sandbox',
    status ENUM('active','revoked','expired') NOT NULL DEFAULT 'active',
    last_used_at TIMESTAMP NULL DEFAULT NULL,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ak_developer (developer_id),
    INDEX idx_ak_status (status),
    INDEX idx_ak_prefix (api_key_prefix)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_applications (
    id VARCHAR(36) PRIMARY KEY,
    developer_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    logo_url VARCHAR(500) DEFAULT NULL,
    redirect_uris JSON DEFAULT NULL,
    webhook_url VARCHAR(500) DEFAULT NULL,
    webhook_secret VARCHAR(255) DEFAULT NULL,
    oauth_client_id VARCHAR(255) DEFAULT NULL,
    oauth_client_secret_hash VARCHAR(255) DEFAULT NULL,
    status ENUM('draft','active','suspended') NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_aapp_developer (developer_id),
    INDEX idx_aapp_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_usage_logs (
    id VARCHAR(36) PRIMARY KEY,
    api_key_id VARCHAR(36) NOT NULL,
    developer_id VARCHAR(36) NOT NULL,
    endpoint VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INT NOT NULL,
    response_time_ms INT NOT NULL DEFAULT 0,
    request_size_bytes INT NOT NULL DEFAULT 0,
    response_size_bytes INT NOT NULL DEFAULT 0,
    ip_address VARCHAR(45) DEFAULT NULL,
    user_agent TEXT DEFAULT NULL,
    error_code VARCHAR(64) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_aul_api_key (api_key_id),
    INDEX idx_aul_developer (developer_id),
    INDEX idx_aul_created (created_at),
    INDEX idx_aul_endpoint (endpoint(100))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_usage_quotas (
    id VARCHAR(36) PRIMARY KEY,
    developer_id VARCHAR(36) NOT NULL,
    period ENUM('daily','monthly') NOT NULL DEFAULT 'daily',
    quota_limit INT NOT NULL DEFAULT 10000,
    quota_used INT NOT NULL DEFAULT 0,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_auq_dev_period (developer_id, period, period_start),
    INDEX idx_auq_period (period_start, period_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS webhook_subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    developer_id VARCHAR(36) NOT NULL,
    api_key_id VARCHAR(36) NOT NULL,
    event_types JSON NOT NULL,
    callback_url VARCHAR(500) NOT NULL,
    secret VARCHAR(255) NOT NULL,
    status ENUM('active','paused','disabled') NOT NULL DEFAULT 'active',
    failure_count INT NOT NULL DEFAULT 0,
    last_triggered_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ws_developer (developer_id),
    INDEX idx_ws_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS webhook_delivery_logs (
    id VARCHAR(36) PRIMARY KEY,
    subscription_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL,
    callback_url VARCHAR(500) NOT NULL,
    http_status INT DEFAULT NULL,
    response_body TEXT DEFAULT NULL,
    response_time_ms INT NOT NULL DEFAULT 0,
    status ENUM('pending','delivered','failed','retrying') NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_wdl_subscription (subscription_id),
    INDEX idx_wdl_status (status),
    INDEX idx_wdl_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
