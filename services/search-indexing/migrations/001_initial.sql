-- Search Indexing Service - Initial Schema
-- Search index management, query analytics, and autocomplete

CREATE TABLE IF NOT EXISTS search_index_tasks (
    id VARCHAR(36) PRIMARY KEY,
    task_type ENUM('full_reindex','incremental','delete','update_mapping') NOT NULL DEFAULT 'incremental',
    index_name VARCHAR(255) NOT NULL,
    status ENUM('pending','queued','running','completed','failed','cancelled') NOT NULL DEFAULT 'pending',
    priority INT NOT NULL DEFAULT 5,
    source_table VARCHAR(100) DEFAULT NULL,
    records_total BIGINT NOT NULL DEFAULT 0,
    records_processed BIGINT NOT NULL DEFAULT 0,
    records_failed BIGINT NOT NULL DEFAULT 0,
    error_message TEXT DEFAULT NULL,
    worker_id VARCHAR(36) DEFAULT NULL,
    started_at TIMESTAMP NULL DEFAULT NULL,
    completed_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_sit_status (status, priority),
    INDEX idx_sit_index (index_name),
    INDEX idx_sit_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS search_index_metadata (
    id VARCHAR(36) PRIMARY KEY,
    index_name VARCHAR(255) NOT NULL UNIQUE,
    index_type ENUM('products','sellers','categories','orders','users','all') NOT NULL,
    document_count BIGINT NOT NULL DEFAULT 0,
    index_size_bytes BIGINT NOT NULL DEFAULT 0,
    shard_count INT NOT NULL DEFAULT 1,
    replica_count INT NOT NULL DEFAULT 1,
    mapping_version INT NOT NULL DEFAULT 1,
    settings_json JSON DEFAULT NULL,
    mapping_json JSON DEFAULT NULL,
    last_reindexed_at TIMESTAMP NULL DEFAULT NULL,
    status ENUM('active','building','error','deprecated') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_sim_type (index_type, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS search_query_logs (
    id VARCHAR(36) PRIMARY KEY,
    query_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) DEFAULT NULL,
    session_id VARCHAR(36) DEFAULT NULL,
    query_text VARCHAR(1000) NOT NULL,
    normalized_query VARCHAR(1000) NOT NULL,
    filters JSON DEFAULT NULL,
    sort_by VARCHAR(64) DEFAULT NULL,
    page INT NOT NULL DEFAULT 1,
    page_size INT NOT NULL DEFAULT 20,
    result_count BIGINT NOT NULL DEFAULT 0,
    response_time_ms INT NOT NULL DEFAULT 0,
    clicked_results JSON DEFAULT NULL,
    has_results BOOLEAN NOT NULL DEFAULT TRUE,
    is_zero_results BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_sql_user (user_id),
    INDEX idx_sql_query (normalized_query(255)),
    INDEX idx_sql_created (created_at),
    INDEX idx_sql_zero (is_zero_results)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS autocomplete_suggestions (
    id VARCHAR(36) PRIMARY KEY,
    suggestion_text VARCHAR(500) NOT NULL,
    normalized_text VARCHAR(500) NOT NULL,
    suggestion_type ENUM('product','category','brand','trending','recent','popular') NOT NULL DEFAULT 'popular',
    score FLOAT NOT NULL DEFAULT 0,
    result_count BIGINT NOT NULL DEFAULT 0,
    click_count BIGINT NOT NULL DEFAULT 0,
    search_count BIGINT NOT NULL DEFAULT 0,
    locale VARCHAR(10) NOT NULL DEFAULT 'en',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    valid_from TIMESTAMP NULL DEFAULT NULL,
    valid_until TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_as_text (normalized_text(255)),
    INDEX idx_as_type (suggestion_type, locale, is_active),
    INDEX idx_as_score (score DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS search_synonyms (
    id VARCHAR(36) PRIMARY KEY,
    term VARCHAR(255) NOT NULL,
    synonyms JSON NOT NULL,
    locale VARCHAR(10) NOT NULL DEFAULT 'en',
    is_bidirectional BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_ss_term_locale (term, locale),
    INDEX idx_ss_locale (locale, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS search_analytics_daily (
    id VARCHAR(36) PRIMARY KEY,
    analytics_date DATE NOT NULL,
    total_queries BIGINT NOT NULL DEFAULT 0,
    unique_queries BIGINT NOT NULL DEFAULT 0,
    zero_result_queries BIGINT NOT NULL DEFAULT 0,
    avg_response_time_ms FLOAT NOT NULL DEFAULT 0,
    p95_response_time_ms FLOAT NOT NULL DEFAULT 0,
    p99_response_time_ms FLOAT NOT NULL DEFAULT 0,
    total_clicks BIGINT NOT NULL DEFAULT 0,
    ctr FLOAT NOT NULL DEFAULT 0,
    top_queries JSON DEFAULT NULL,
    top_zero_result_queries JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_sad_date (analytics_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
