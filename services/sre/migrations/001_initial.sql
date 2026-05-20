-- SRE Service - Initial Schema
-- Incident management, SLO tracking, runbooks, and on-call scheduling

CREATE TABLE IF NOT EXISTS slo_definitions (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    slo_type ENUM('availability','latency','error_rate','throughput') NOT NULL,
    target_percentage FLOAT NOT NULL DEFAULT 99.9,
    latency_threshold_ms INT DEFAULT NULL,
    latency_percentile FLOAT DEFAULT NULL,
    window_days INT NOT NULL DEFAULT 30,
    burn_rate_alert_threshold FLOAT NOT NULL DEFAULT 14.4,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_slo_svc_name (service_name, name),
    INDEX idx_slo_service (service_name, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS slo_measurements (
    id VARCHAR(36) PRIMARY KEY,
    slo_id VARCHAR(36) NOT NULL,
    measurement_date DATE NOT NULL,
    good_events BIGINT NOT NULL DEFAULT 0,
    total_events BIGINT NOT NULL DEFAULT 0,
    slo_percentage FLOAT NOT NULL DEFAULT 100.0,
    error_budget_remaining FLOAT NOT NULL DEFAULT 100.0,
    burn_rate FLOAT NOT NULL DEFAULT 0,
    is_violation BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_sm_slo_date (slo_id, measurement_date),
    INDEX idx_sm_date (measurement_date),
    INDEX idx_sm_violation (is_violation)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS incidents (
    id VARCHAR(36) PRIMARY KEY,
    incident_number VARCHAR(36) NOT NULL UNIQUE,
    title VARCHAR(500) NOT NULL,
    description TEXT DEFAULT NULL,
    severity ENUM('P0','P1','P2','P3','P4') NOT NULL DEFAULT 'P3',
    status ENUM('detected','investigating','mitigating','monitoring','resolved','postmortem') NOT NULL DEFAULT 'detected',
    service_name VARCHAR(100) NOT NULL,
    affected_services JSON DEFAULT NULL,
    root_cause TEXT DEFAULT NULL,
    resolution TEXT DEFAULT NULL,
    detected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    acknowledged_at TIMESTAMP NULL DEFAULT NULL,
    mitigated_at TIMESTAMP NULL DEFAULT NULL,
    resolved_at TIMESTAMP NULL DEFAULT NULL,
    postmortem_due_at TIMESTAMP NULL DEFAULT NULL,
    postmortem_url VARCHAR(500) DEFAULT NULL,
    commander_id VARCHAR(36) DEFAULT NULL,
    communication_channel VARCHAR(255) DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_inc_severity (severity, status),
    INDEX idx_inc_service (service_name),
    INDEX idx_inc_status (status),
    INDEX idx_inc_detected (detected_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS incident_timeline (
    id VARCHAR(36) PRIMARY KEY,
    incident_id VARCHAR(36) NOT NULL,
    event_type ENUM('detected','acknowledged','escalated','mitigated','resolved','note','action','communication') NOT NULL,
    description TEXT NOT NULL,
    actor_id VARCHAR(36) DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_it_incident (incident_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS on_call_schedules (
    id VARCHAR(36) PRIMARY KEY,
    team_name VARCHAR(100) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    rotation_type ENUM('weekly','daily','custom') NOT NULL DEFAULT 'weekly',
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT TRUE,
    escalation_order INT NOT NULL DEFAULT 1,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ocs_team (team_name),
    INDEX idx_ocs_user (user_id),
    INDEX idx_ocs_time (start_time, end_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS runbooks (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    category ENUM('incident_response','deployment','monitoring','backup','security','maintenance') NOT NULL DEFAULT 'incident_response',
    content_md TEXT NOT NULL,
    version INT NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    tags JSON DEFAULT NULL,
    created_by VARCHAR(36) DEFAULT NULL,
    last_reviewed_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_rb_service (service_name, category),
    INDEX idx_rb_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS alert_rules (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    condition_type ENUM('threshold','anomaly','absence','rate') NOT NULL DEFAULT 'threshold',
    operator ENUM('gt','gte','lt','lte','eq','neq') NOT NULL,
    threshold_value FLOAT NOT NULL,
    duration_sec INT NOT NULL DEFAULT 300,
    severity ENUM('critical','warning','info') NOT NULL DEFAULT 'warning',
    notification_channels JSON NOT NULL,
    runbook_id VARCHAR(36) DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ar_service (service_name, is_active),
    INDEX idx_ar_severity (severity)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
