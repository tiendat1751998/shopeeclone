-- Service Mesh Service - Initial Schema
-- Service discovery, traffic management, and mesh observability

CREATE TABLE IF NOT EXISTS mesh_services (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL UNIQUE,
    namespace VARCHAR(100) NOT NULL DEFAULT 'default',
    version VARCHAR(20) NOT NULL DEFAULT 'v1',
    protocol ENUM('http','grpc','tcp','websocket') NOT NULL DEFAULT 'http',
    port INT NOT NULL DEFAULT 8080,
    health_check_path VARCHAR(255) NOT NULL DEFAULT '/health',
    health_check_interval_sec INT NOT NULL DEFAULT 10,
    status ENUM('healthy','degraded','unhealthy','offline') NOT NULL DEFAULT 'healthy',
    region VARCHAR(36) DEFAULT NULL,
    labels JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ms_namespace (namespace),
    INDEX idx_ms_status (status),
    INDEX idx_ms_region (region)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mesh_endpoints (
    id VARCHAR(36) PRIMARY KEY,
    service_id VARCHAR(36) NOT NULL,
    pod_name VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    port INT NOT NULL DEFAULT 8080,
    weight INT NOT NULL DEFAULT 100,
    status ENUM('ready','not_ready','terminating','starting') NOT NULL DEFAULT 'starting',
    last_health_check TIMESTAMP NULL DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_me_service (service_id),
    INDEX idx_me_status (status),
    INDEX idx_me_ip (ip_address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS traffic_policies (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    policy_type ENUM('canary','blue_green','ab_test','circuit_breaker','retry','timeout','rate_limit','fault_injection') NOT NULL,
    config_json JSON NOT NULL,
    priority INT NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    effective_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    effective_until TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tp_service (service_name),
    INDEX idx_tp_type (policy_type, is_active),
    INDEX idx_tp_active (is_active, priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS traffic_splits (
    id VARCHAR(36) PRIMARY KEY,
    policy_id VARCHAR(36) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    destination_service VARCHAR(100) NOT NULL,
    destination_version VARCHAR(20) NOT NULL,
    weight INT NOT NULL DEFAULT 0,
    match_headers JSON DEFAULT NULL,
    match_query_params JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ts_policy (policy_id),
    INDEX idx_ts_service (service_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mesh_certificates (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    certificate_type ENUM('tls','mtls','jwt') NOT NULL DEFAULT 'tls',
    certificate_pem TEXT NOT NULL,
    private_key_encrypted TEXT NOT NULL,
    ca_pem TEXT DEFAULT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    auto_rotate BOOLEAN NOT NULL DEFAULT TRUE,
    status ENUM('active','expiring','expired','rotating') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_mc_service (service_name),
    INDEX idx_mc_status (status),
    INDEX idx_mc_expiry (valid_until)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mesh_access_policies (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    source_service VARCHAR(100) NOT NULL,
    destination_service VARCHAR(100) NOT NULL,
    allowed_methods JSON DEFAULT NULL,
    allowed_paths JSON DEFAULT NULL,
    action ENUM('allow','deny','audit','require_auth') NOT NULL DEFAULT 'allow',
    priority INT NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_map_source (source_service),
    INDEX idx_map_dest (destination_service),
    INDEX idx_map_active (is_active, priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS mesh_metrics_snapshot (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    metric_time TIMESTAMP NOT NULL,
    request_count BIGINT NOT NULL DEFAULT 0,
    error_count BIGINT NOT NULL DEFAULT 0,
    p50_latency_ms FLOAT NOT NULL DEFAULT 0,
    p95_latency_ms FLOAT NOT NULL DEFAULT 0,
    p99_latency_ms FLOAT NOT NULL DEFAULT 0,
    requests_per_second FLOAT NOT NULL DEFAULT 0,
    error_rate FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_mms_svc_time (service_name, metric_time),
    INDEX idx_mms_time (metric_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
