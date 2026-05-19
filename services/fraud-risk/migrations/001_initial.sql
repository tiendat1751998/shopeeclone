-- Fraud Risk Service - Initial Schema
-- Risk scoring, fraud rules, device fingerprinting, and case management

CREATE TABLE IF NOT EXISTS risk_rules (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    rule_type ENUM('velocity','amount','geolocation','device','behavior','composite','ml_model') NOT NULL,
    category ENUM('payment','account','order','listing','promotion','shipping') NOT NULL DEFAULT 'payment',
    conditions JSON NOT NULL,
    risk_score INT NOT NULL DEFAULT 0,
    action ENUM('allow','flag','review','block','challenge','escalate') NOT NULL DEFAULT 'flag',
    priority INT NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    effective_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    effective_until TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_rr_category (category, is_active),
    INDEX idx_rr_type (rule_type, is_active),
    INDEX idx_rr_priority (priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS risk_assessments (
    id VARCHAR(36) PRIMARY KEY,
    entity_type ENUM('payment','order','account','listing','user','session') NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) DEFAULT NULL,
    risk_score INT NOT NULL DEFAULT 0,
    risk_level ENUM('low','medium','high','critical') NOT NULL DEFAULT 'low',
    decision ENUM('allow','flag','review','block','challenge') NOT NULL DEFAULT 'allow',
    triggered_rules JSON DEFAULT NULL,
    signals JSON DEFAULT NULL,
    device_fingerprint VARCHAR(255) DEFAULT NULL,
    ip_address VARCHAR(45) DEFAULT NULL,
    geo_country VARCHAR(3) DEFAULT NULL,
    assessed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ra_entity (entity_type, entity_id),
    INDEX idx_ra_user (user_id),
    INDEX idx_ra_decision (decision),
    INDEX idx_ra_score (risk_score DESC),
    INDEX idx_ra_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS device_fingerprints (
    id VARCHAR(36) PRIMARY KEY,
    fingerprint_hash VARCHAR(255) NOT NULL UNIQUE,
    user_id VARCHAR(36) DEFAULT NULL,
    device_type VARCHAR(64) DEFAULT NULL,
    os VARCHAR(64) DEFAULT NULL,
    browser VARCHAR(64) DEFAULT NULL,
    screen_resolution VARCHAR(20) DEFAULT NULL,
    timezone VARCHAR(50) DEFAULT NULL,
    language VARCHAR(10) DEFAULT NULL,
    ip_address VARCHAR(45) DEFAULT NULL,
    first_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    trust_score FLOAT NOT NULL DEFAULT 50.0,
    is_blacklisted BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSON DEFAULT NULL,
    INDEX idx_df_user (user_id),
    INDEX idx_df_trust (trust_score),
    INDEX idx_df_blacklisted (is_blacklisted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS fraud_cases (
    id VARCHAR(36) PRIMARY KEY,
    case_number VARCHAR(36) NOT NULL UNIQUE,
    entity_type ENUM('payment','order','account','listing','user') NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) DEFAULT NULL,
    risk_assessment_id VARCHAR(36) DEFAULT NULL,
    status ENUM('open','under_review','confirmed_fraud','false_positive','escalated','closed') NOT NULL DEFAULT 'open',
    priority ENUM('low','medium','high','critical') NOT NULL DEFAULT 'medium',
    fraud_type ENUM('payment_fraud','account_takeover','promotion_abuse','seller_fraud','friendly_fraud','synthetic_identity','other') DEFAULT NULL,
    amount_at_risk BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    evidence JSON DEFAULT NULL,
    notes TEXT DEFAULT NULL,
    assigned_to VARCHAR(36) DEFAULT NULL,
    resolved_at TIMESTAMP NULL DEFAULT NULL,
    resolution_notes TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_fc_status (status),
    INDEX idx_fc_user (user_id),
    INDEX idx_fc_entity (entity_type, entity_id),
    INDEX idx_fc_assigned (assigned_to, status),
    INDEX idx_fc_priority (priority, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS velocity_counters (
    id VARCHAR(36) PRIMARY KEY,
    entity_type ENUM('user','device','ip','card','account') NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    window_type ENUM('minute','hour','day','week','month') NOT NULL,
    window_start TIMESTAMP NOT NULL,
    window_end TIMESTAMP NOT NULL,
    transaction_count INT NOT NULL DEFAULT 0,
    total_amount BIGINT NOT NULL DEFAULT 0,
    unique_recipients INT NOT NULL DEFAULT 0,
    unique_ips INT NOT NULL DEFAULT 0,
    unique_devices INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_vc_entity_window (entity_type, entity_id, window_type, window_start),
    INDEX idx_vc_window (window_start, window_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS blacklist_entries (
    id VARCHAR(36) PRIMARY KEY,
    entry_type ENUM('ip','email','phone','device','card_bin','card_hash','user','address') NOT NULL,
    entry_value VARCHAR(255) NOT NULL,
    reason TEXT NOT NULL,
    source ENUM('manual','rule','ml_model','external_feed','report') NOT NULL DEFAULT 'manual',
    severity ENUM('low','medium','high','critical') NOT NULL DEFAULT 'medium',
    expires_at TIMESTAMP NULL DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by VARCHAR(36) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_be_type_value (entry_type, entry_value),
    INDEX idx_be_active (is_active, entry_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
