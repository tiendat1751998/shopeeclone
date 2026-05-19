CREATE TABLE IF NOT EXISTS index_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    idempotency_key VARCHAR(255) UNIQUE,
    error TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_index_tasks_status ON index_tasks(status);
CREATE INDEX idx_index_tasks_idempotency_key ON index_tasks(idempotency_key);
CREATE INDEX idx_index_tasks_created_at ON index_tasks(created_at);
