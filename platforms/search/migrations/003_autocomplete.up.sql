CREATE TABLE IF NOT EXISTS autocomplete_suggestions (
    id BIGSERIAL PRIMARY KEY,
    prefix VARCHAR(100) NOT NULL,
    suggestion TEXT NOT NULL,
    score DOUBLE PRECISION NOT NULL DEFAULT 0,
    type VARCHAR(50) NOT NULL DEFAULT 'product',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_autocomplete_prefix ON autocomplete_suggestions(prefix);
CREATE INDEX idx_autocomplete_score ON autocomplete_suggestions(score DESC);

CREATE TABLE IF NOT EXISTS trending_queries (
    id BIGSERIAL PRIMARY KEY,
    query VARCHAR(500) NOT NULL UNIQUE,
    score DOUBLE PRECISION NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trending_queries_score ON trending_queries(score DESC);
