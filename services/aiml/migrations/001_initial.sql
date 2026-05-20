-- AI/ML Service - Initial Schema
-- Model registry, feature store, predictions, and A/B testing

CREATE TABLE IF NOT EXISTS ml_models (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    model_type ENUM('classification','regression','ranking','recommendation','nlp','cv','embedding','anomaly') NOT NULL,
    framework VARCHAR(64) NOT NULL,
    artifact_path VARCHAR(500) NOT NULL,
    artifact_size_bytes BIGINT NOT NULL DEFAULT 0,
    input_schema JSON NOT NULL,
    output_schema JSON NOT NULL,
    hyperparameters JSON DEFAULT NULL,
    metrics JSON DEFAULT NULL,
    status ENUM('training','staging','production','deprecated','failed') NOT NULL DEFAULT 'training',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT DEFAULT NULL,
    created_by VARCHAR(36) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_mm_name_ver (name, version),
    INDEX idx_mm_type (model_type, status),
    INDEX idx_mm_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS model_deployments (
    id VARCHAR(36) PRIMARY KEY,
    model_id VARCHAR(36) NOT NULL,
    deployment_name VARCHAR(255) NOT NULL,
    environment ENUM('staging','production','canary') NOT NULL DEFAULT 'staging',
    endpoint_url VARCHAR(500) DEFAULT NULL,
    replicas INT NOT NULL DEFAULT 1,
    cpu_request VARCHAR(20) DEFAULT '500m',
    memory_request VARCHAR(20) DEFAULT '1Gi',
    gpu_request INT NOT NULL DEFAULT 0,
    status ENUM('deploying','running','stopped','failed','scaling') NOT NULL DEFAULT 'deploying',
    health_check_url VARCHAR(500) DEFAULT NULL,
    deployed_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_md_model (model_id),
    INDEX idx_md_env (environment, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS feature_groups (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT DEFAULT NULL,
    entity_type ENUM('user','product','order','session','seller','category') NOT NULL,
    storage_type ENUM('online','offline','both') NOT NULL DEFAULT 'both',
    ttl_seconds INT DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_fg_entity (entity_type, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS features (
    id VARCHAR(36) PRIMARY KEY,
    group_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    data_type ENUM('int','float','string','boolean','array','json','timestamp') NOT NULL,
    default_value VARCHAR(255) DEFAULT NULL,
    transformation TEXT DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_f_group_name (group_id, name),
    INDEX idx_f_group (group_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS prediction_logs (
    id VARCHAR(36) PRIMARY KEY,
    model_id VARCHAR(36) NOT NULL,
    deployment_id VARCHAR(36) DEFAULT NULL,
    request_id VARCHAR(36) NOT NULL,
    entity_type VARCHAR(64) DEFAULT NULL,
    entity_id VARCHAR(36) DEFAULT NULL,
    input_features JSON NOT NULL,
    output_prediction JSON NOT NULL,
    confidence_score FLOAT DEFAULT NULL,
    latency_ms INT NOT NULL DEFAULT 0,
    status ENUM('success','error','timeout') NOT NULL DEFAULT 'success',
    error_message TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_pl_model (model_id),
    INDEX idx_pl_request (request_id),
    INDEX idx_pl_entity (entity_type, entity_id),
    INDEX idx_pl_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ab_experiments (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT DEFAULT NULL,
    hypothesis TEXT DEFAULT NULL,
    status ENUM('draft','running','paused','completed','cancelled') NOT NULL DEFAULT 'draft',
    model_a_id VARCHAR(36) NOT NULL,
    model_b_id VARCHAR(36) NOT NULL,
    traffic_split_a FLOAT NOT NULL DEFAULT 50.0,
    traffic_split_b FLOAT NOT NULL DEFAULT 50.0,
    target_metric VARCHAR(100) NOT NULL,
    min_sample_size INT NOT NULL DEFAULT 10000,
    confidence_level FLOAT NOT NULL DEFAULT 95.0,
    started_at TIMESTAMP NULL DEFAULT NULL,
    ended_at TIMESTAMP NULL DEFAULT NULL,
    winner ENUM('a','b','tie','inconclusive') DEFAULT NULL,
    results JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ae_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS training_jobs (
    id VARCHAR(36) PRIMARY KEY,
    model_id VARCHAR(36) NOT NULL,
    job_name VARCHAR(255) NOT NULL,
    dataset_path VARCHAR(500) NOT NULL,
    hyperparameters JSON NOT NULL,
    status ENUM('queued','running','completed','failed','cancelled') NOT NULL DEFAULT 'queued',
    progress_percent FLOAT NOT NULL DEFAULT 0,
    metrics JSON DEFAULT NULL,
    error_message TEXT DEFAULT NULL,
    started_at TIMESTAMP NULL DEFAULT NULL,
    completed_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tj_model (model_id),
    INDEX idx_tj_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
