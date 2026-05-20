-- Gateway Service - Initial Schema
-- API Gateway routing, rate limiting, and session management

CREATE TABLE IF NOT EXISTS api_routes (
    id VARCHAR(36) PRIMARY KEY,
    path VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    strip_path BOOLEAN NOT NULL DEFAULT FALSE,
    preserve_host BOOLEAN NOT NULL DEFAULT FALSE,
    timeout_ms INT NOT NULL DEFAULT 30000,
    retry_count INT NOT NULL DEFAULT 3,
    circuit_breaker_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    rate_limit_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    rate_limit_per_minute INT NOT NULL DEFAULT 1000,
    auth_required BOOLEAN NOT NULL DEFAULT TRUE,
    cors_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_ar_path_method (path, method),
    INDEX idx_ar_service (service_name, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS rate_limit_rules (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    scope ENUM('global','service','route','ip','user') NOT NULL DEFAULT 'global',
    target VARCHAR(255) DEFAULT NULL,
    requests_per_minute INT NOT NULL DEFAULT 100,
    burst_size INT NOT NULL DEFAULT 150,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_rlr_scope (scope, target)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS gateway_sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) DEFAULT NULL,
    session_token_hash VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(45) DEFAULT NULL,
    user_agent TEXT DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_gs_user (user_id, is_active),
    INDEX idx_gs_token (session_token_hash),
    INDEX idx_gs_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS service_health_checks (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    endpoint VARCHAR(500) NOT NULL,
    check_interval_sec INT NOT NULL DEFAULT 30,
    timeout_sec INT NOT NULL DEFAULT 5,
    healthy_threshold INT NOT NULL DEFAULT 2,
    unhealthy_threshold INT NOT NULL DEFAULT 3,
    last_check_at TIMESTAMP NULL DEFAULT NULL,
    last_status ENUM('healthy','unhealthy','unknown') NOT NULL DEFAULT 'unknown',
    consecutive_failures INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_shc_service (service_name),
    INDEX idx_shc_status (last_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS circuit_breaker_states (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL UNIQUE,
    state ENUM('closed','open','half_open') NOT NULL DEFAULT 'closed',
    failure_count INT NOT NULL DEFAULT 0,
    success_count INT NOT NULL DEFAULT 0,
    last_failure_at TIMESTAMP NULL DEFAULT NULL,
    last_success_at TIMESTAMP NULL DEFAULT NULL,
    opened_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_cbs_service (service_name),
    INDEX idx_cbs_state (state)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
