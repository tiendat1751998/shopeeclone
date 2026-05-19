-- API Gateway Service - Initial Schema
-- API versioning, consumer management, and request/response transformation

CREATE TABLE IF NOT EXISTS api_consumers (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    organization VARCHAR(255) DEFAULT NULL,
    contact_person VARCHAR(255) DEFAULT NULL,
    status ENUM('active','inactive','suspended','pending_approval') NOT NULL DEFAULT 'pending_approval',
    tier ENUM('free','standard','premium','enterprise') NOT NULL DEFAULT 'free',
    monthly_quota BIGINT NOT NULL DEFAULT 10000,
    current_month_usage BIGINT NOT NULL DEFAULT 0,
    quota_reset_day INT NOT NULL DEFAULT 1,
    billing_email VARCHAR(255) DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ac_status (status),
    INDEX idx_ac_tier (tier),
    INDEX idx_ac_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_consumer_keys (
    id VARCHAR(36) PRIMARY KEY,
    consumer_id VARCHAR(36) NOT NULL,
    key_type ENUM('api_key','oauth2','jwt','mtls') NOT NULL DEFAULT 'api_key',
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    key_prefix VARCHAR(8) NOT NULL,
    scopes JSON NOT NULL,
    rate_limit_per_second INT NOT NULL DEFAULT 10,
    rate_limit_per_minute INT NOT NULL DEFAULT 600,
    rate_limit_per_day INT NOT NULL DEFAULT 50000,
    allowed_ips JSON DEFAULT NULL,
    allowed_origins JSON DEFAULT NULL,
    environment ENUM('sandbox','production','both') NOT NULL DEFAULT 'sandbox',
    status ENUM('active','revoked','expired') NOT NULL DEFAULT 'active',
    last_used_at TIMESTAMP NULL DEFAULT NULL,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ack_consumer (consumer_id),
    INDEX idx_ack_status (status),
    INDEX idx_ack_prefix (key_prefix)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_versions (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    version VARCHAR(20) NOT NULL,
    base_path VARCHAR(255) NOT NULL,
    status ENUM('draft','active','deprecated','retired') NOT NULL DEFAULT 'draft',
    changelog TEXT DEFAULT NULL,
    deprecated_at TIMESTAMP NULL DEFAULT NULL,
    retired_at TIMESTAMP NULL DEFAULT NULL,
    sunset_date TIMESTAMP NULL DEFAULT NULL,
    documentation_url VARCHAR(500) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_av_svc_ver (service_name, version),
    INDEX idx_av_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_request_logs (
    id VARCHAR(36) PRIMARY KEY,
    request_id VARCHAR(36) NOT NULL,
    consumer_id VARCHAR(36) DEFAULT NULL,
    api_key_id VARCHAR(36) DEFAULT NULL,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    query_params TEXT DEFAULT NULL,
    request_headers JSON DEFAULT NULL,
    request_body_size INT NOT NULL DEFAULT 0,
    response_status INT NOT NULL,
    response_headers JSON DEFAULT NULL,
    response_body_size INT NOT NULL DEFAULT 0,
    response_time_ms INT NOT NULL DEFAULT 0,
    ip_address VARCHAR(45) DEFAULT NULL,
    user_agent TEXT DEFAULT NULL,
    error_code VARCHAR(64) DEFAULT NULL,
    error_message TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_arl_consumer (consumer_id),
    INDEX idx_arl_request (request_id),
    INDEX idx_arl_path (path(100)),
    INDEX idx_arl_status (response_status),
    INDEX idx_arl_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_transformations (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    version_id VARCHAR(36) DEFAULT NULL,
    path_pattern VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL DEFAULT '*',
    transform_type ENUM('request','response','both') NOT NULL DEFAULT 'both',
    headers_add JSON DEFAULT NULL,
    headers_remove JSON DEFAULT NULL,
    body_template TEXT DEFAULT NULL,
    query_add JSON DEFAULT NULL,
    query_remove JSON DEFAULT NULL,
    priority INT NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_at_service (service_name),
    INDEX idx_at_version (version_id),
    INDEX idx_at_active (is_active, priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS api_health_metrics (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    metric_window_start TIMESTAMP NOT NULL,
    metric_window_end TIMESTAMP NOT NULL,
    total_requests BIGINT NOT NULL DEFAULT 0,
    successful_requests BIGINT NOT NULL DEFAULT 0,
    error_requests BIGINT NOT NULL DEFAULT 0,
    avg_response_time_ms FLOAT NOT NULL DEFAULT 0,
    p50_response_time_ms FLOAT NOT NULL DEFAULT 0,
    p95_response_time_ms FLOAT NOT NULL DEFAULT 0,
    p99_response_time_ms FLOAT NOT NULL DEFAULT 0,
    requests_per_second FLOAT NOT NULL DEFAULT 0,
    unique_consumers INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_ahm_svc_window (service_name, metric_window_start),
    INDEX idx_ahm_window (metric_window_start, metric_window_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
