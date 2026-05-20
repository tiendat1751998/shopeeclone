-- Global Infrastructure Service - Initial Schema
-- Multi-region config, DNS management, and infrastructure inventory

CREATE TABLE IF NOT EXISTS regions (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    provider ENUM('aws','gcp','azure','on_premise','hybrid') NOT NULL,
    country_code VARCHAR(3) NOT NULL,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_r_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS infrastructure_nodes (
    id VARCHAR(36) PRIMARY KEY,
    region_id VARCHAR(36) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    node_type ENUM('kubernetes_node','database','cache','queue','load_balancer','cdn','storage','monitoring','bastion') NOT NULL,
    role ENUM('primary','secondary','replica','worker','master','edge') NOT NULL DEFAULT 'worker',
    cpu_cores INT NOT NULL DEFAULT 0,
    memory_gb INT NOT NULL DEFAULT 0,
    storage_gb INT NOT NULL DEFAULT 0,
    os VARCHAR(64) DEFAULT NULL,
    status ENUM('provisioning','running','maintenance','decommissioned','failed') NOT NULL DEFAULT 'provisioning',
    labels JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_in_region (region_id),
    INDEX idx_in_type (node_type, status),
    INDEX idx_in_ip (ip_address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS dns_records (
    id VARCHAR(36) PRIMARY KEY,
    zone VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    record_type ENUM('A','AAAA','CNAME','MX','TXT','SRV','NS','CAA') NOT NULL,
    value VARCHAR(500) NOT NULL,
    ttl INT NOT NULL DEFAULT 300,
    priority INT DEFAULT NULL,
    region_id VARCHAR(36) DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_dr_zone_name_type (zone, name, record_type, value(100)),
    INDEX idx_dr_zone (zone),
    INDEX idx_dr_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ssl_certificates (
    id VARCHAR(36) PRIMARY KEY,
    domain VARCHAR(255) NOT NULL,
    issuer VARCHAR(255) NOT NULL,
    certificate_pem TEXT NOT NULL,
    private_key_encrypted TEXT NOT NULL,
    chain_pem TEXT DEFAULT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    auto_renew BOOLEAN NOT NULL DEFAULT TRUE,
    status ENUM('active','expiring','expired','revoked','pending') NOT NULL DEFAULT 'active',
    last_renewed_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_sc_domain (domain),
    INDEX idx_sc_status (status),
    INDEX idx_sc_expiry (valid_until)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS infrastructure_changes (
    id VARCHAR(36) PRIMARY KEY,
    change_type ENUM('provision','update','decommission','scale','config_change','security_patch') NOT NULL,
    resource_type VARCHAR(64) NOT NULL,
    resource_id VARCHAR(36) NOT NULL,
    region_id VARCHAR(36) DEFAULT NULL,
    description TEXT NOT NULL,
    status ENUM('pending','approved','in_progress','completed','failed','rolled_back') NOT NULL DEFAULT 'pending',
    requested_by VARCHAR(36) NOT NULL,
    approved_by VARCHAR(36) DEFAULT NULL,
    approved_at TIMESTAMP NULL DEFAULT NULL,
    started_at TIMESTAMP NULL DEFAULT NULL,
    completed_at TIMESTAMP NULL DEFAULT NULL,
    rollback_plan TEXT DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ic_status (status),
    INDEX idx_ic_resource (resource_type, resource_id),
    INDEX idx_ic_region (region_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
