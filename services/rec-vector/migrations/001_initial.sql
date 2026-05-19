-- Recommendation Vector Service - Initial Schema
-- Vector embeddings, similarity search, and recommendation models

CREATE TABLE IF NOT EXISTS vector_collections (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT DEFAULT NULL,
    entity_type ENUM('product','user','seller','category','query','image') NOT NULL,
    vector_dimension INT NOT NULL DEFAULT 768,
    distance_metric ENUM('cosine','euclidean','dot','manhattan') NOT NULL DEFAULT 'cosine',
    index_type ENUM('flat','ivf','hnsw','pq') NOT NULL DEFAULT 'hnsw',
    index_params JSON DEFAULT NULL,
    total_vectors BIGINT NOT NULL DEFAULT 0,
    status ENUM('active','building','reindexing','error','deprecated') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_vc_entity (entity_type, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS vector_metadata (
    id VARCHAR(36) PRIMARY KEY,
    collection_id VARCHAR(36) NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    vector_id VARCHAR(36) NOT NULL,
    metadata_json JSON DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_vm_coll_entity (collection_id, entity_id),
    INDEX idx_vm_collection (collection_id, is_active),
    INDEX idx_vm_vector (vector_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS embedding_models (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    model_type ENUM('text','image','multimodal','graph','sequence') NOT NULL,
    framework VARCHAR(64) NOT NULL,
    vector_dimension INT NOT NULL DEFAULT 768,
    artifact_path VARCHAR(500) NOT NULL,
    preprocessing_config JSON DEFAULT NULL,
    status ENUM('training','staging','production','deprecated') NOT NULL DEFAULT 'training',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    metrics JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_em_name_ver (name, version),
    INDEX idx_em_type (model_type, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS recommendation_results (
    id VARCHAR(36) PRIMARY KEY,
    request_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) DEFAULT NULL,
    session_id VARCHAR(36) DEFAULT NULL,
    collection_id VARCHAR(36) NOT NULL,
    source_entity_type VARCHAR(64) DEFAULT NULL,
    source_entity_id VARCHAR(36) DEFAULT NULL,
    algorithm ENUM('collaborative_filtering','content_based','hybrid','trending','personalized','similar_items','frequently_bought') NOT NULL,
    recommendations JSON NOT NULL,
    latency_ms INT NOT NULL DEFAULT 0,
    context JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_rr_user (user_id),
    INDEX idx_rr_request (request_id),
    INDEX idx_rr_collection (collection_id),
    INDEX idx_rr_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_embedding_profiles (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    collection_id VARCHAR(36) NOT NULL,
    vector_id VARCHAR(36) NOT NULL,
    preference_categories JSON DEFAULT NULL,
    preference_brands JSON DEFAULT NULL,
    price_range_min FLOAT DEFAULT NULL,
    price_range_max FLOAT DEFAULT NULL,
    last_interaction_at TIMESTAMP NULL DEFAULT NULL,
    interaction_count INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_uep_collection (collection_id),
    INDEX idx_uep_vector (vector_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS vector_index_jobs (
    id VARCHAR(36) PRIMARY KEY,
    collection_id VARCHAR(36) NOT NULL,
    job_type ENUM('full_build','incremental_update','optimize','delete') NOT NULL DEFAULT 'incremental_update',
    status ENUM('queued','running','completed','failed','cancelled') NOT NULL DEFAULT 'queued',
    priority INT NOT NULL DEFAULT 5,
    vectors_total BIGINT NOT NULL DEFAULT 0,
    vectors_processed BIGINT NOT NULL DEFAULT 0,
    error_message TEXT DEFAULT NULL,
    worker_id VARCHAR(36) DEFAULT NULL,
    started_at TIMESTAMP NULL DEFAULT NULL,
    completed_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_vij_collection (collection_id),
    INDEX idx_vij_status (status, priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
