-- Production Dashboard Migration 001

CREATE TABLE IF NOT EXISTS service_health (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL UNIQUE,
    status ENUM('healthy', 'degraded', 'unhealthy', 'unknown') NOT NULL DEFAULT 'unknown',
    health_url VARCHAR(500) NOT NULL,
    response_time_ms INT NOT NULL DEFAULT 0,
    last_checked_at TIMESTAMP NULL DEFAULT NULL,
    last_error TEXT DEFAULT NULL,
    version VARCHAR(100) DEFAULT NULL,
    environment VARCHAR(50) NOT NULL DEFAULT 'production',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_service_health_status (status),
    INDEX idx_service_health_env (environment)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS deployments (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    version VARCHAR(100) NOT NULL,
    environment VARCHAR(50) NOT NULL,
    status ENUM('pending', 'in_progress', 'succeeded', 'failed', 'rolled_back') NOT NULL DEFAULT 'pending',
    deployed_by VARCHAR(255) NOT NULL,
    image VARCHAR(500) NOT NULL,
    replicas INT NOT NULL DEFAULT 1,
    ready_replicas INT NOT NULL DEFAULT 0,
    strategy VARCHAR(50) NOT NULL DEFAULT 'rolling',
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP NULL DEFAULT NULL,
    rollback_of VARCHAR(36) DEFAULT NULL,
    notes TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_deployments_service (service_name),
    INDEX idx_deployments_status (status),
    INDEX idx_deployments_env (environment),
    INDEX idx_deployments_started (started_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS incidents (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    severity ENUM('critical', 'high', 'medium', 'low') NOT NULL,
    status ENUM('open', 'acknowledged', 'investigating', 'mitigating', 'resolved', 'closed') NOT NULL DEFAULT 'open',
    service_names TEXT NOT NULL,
    detected_at TIMESTAMP NOT NULL,
    acknowledged_at TIMESTAMP NULL DEFAULT NULL,
    resolved_at TIMESTAMP NULL DEFAULT NULL,
    root_cause TEXT DEFAULT NULL,
    resolution TEXT DEFAULT NULL,
    created_by VARCHAR(255) NOT NULL,
    assigned_to VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_incidents_status (status),
    INDEX idx_incidents_severity (severity),
    INDEX idx_incidents_detected (detected_at),
    INDEX idx_incidents_assigned (assigned_to)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS alert_rules (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    service_name VARCHAR(255) NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    `condition` VARCHAR(10) NOT NULL,
    threshold DOUBLE NOT NULL,
    duration VARCHAR(20) NOT NULL DEFAULT '5m',
    severity ENUM('critical', 'high', 'medium', 'low') NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notify_channels TEXT DEFAULT NULL,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_alert_rules_service (service_name),
    INDEX idx_alert_rules_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    actor VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id VARCHAR(36) DEFAULT NULL,
    details TEXT DEFAULT NULL,
    ip_address VARCHAR(45) DEFAULT NULL,
    user_agent VARCHAR(500) DEFAULT NULL,
    created_at VARCHAR(30) NOT NULL,
    INDEX idx_audit_actor (actor),
    INDEX idx_audit_resource (resource, resource_id),
    INDEX idx_audit_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS service_dependencies (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    depends_on VARCHAR(255) NOT NULL,
    dependency_type ENUM('sync', 'async', 'data') NOT NULL DEFAULT 'sync',
    critical BOOLEAN NOT NULL DEFAULT FALSE,
    created_at VARCHAR(30) NOT NULL,
    INDEX idx_dep_service (service_name),
    INDEX idx_dep_depends_on (depends_on),
    UNIQUE KEY uk_service_dep (service_name, depends_on)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS capacity_metrics (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    current_value DOUBLE NOT NULL,
    max_value DOUBLE NOT NULL,
    unit VARCHAR(20) NOT NULL,
    recorded_at VARCHAR(30) NOT NULL,
    INDEX idx_capacity_service (service_name),
    INDEX idx_capacity_resource (service_name, resource_type),
    INDEX idx_capacity_recorded (recorded_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
