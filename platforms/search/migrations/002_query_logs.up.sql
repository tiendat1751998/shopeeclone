CREATE TABLE IF NOT EXISTS query_logs (
    id BIGSERIAL PRIMARY KEY,
    query VARCHAR(500) NOT NULL,
    normalized_query VARCHAR(500),
    result_count BIGINT NOT NULL DEFAULT 0,
    took_ms BIGINT NOT NULL DEFAULT 0,
    user_id VARCHAR(255),
    session_id VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_query_logs_query ON query_logs(query);
CREATE INDEX idx_query_logs_created_at ON query_logs(created_at);
CREATE INDEX idx_query_logs_normalized ON query_logs(normalized_query);
