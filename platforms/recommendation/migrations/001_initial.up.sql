CREATE TABLE IF NOT EXISTS recommendation_events (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    session_id VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_profiles (
    user_id VARCHAR(255) PRIMARY KEY,
    category_weights JSONB NOT NULL DEFAULT '{}',
    preferred_brands JSONB NOT NULL DEFAULT '{}',
    interest_vector JSONB NOT NULL DEFAULT '{}',
    price_range_min NUMERIC(12,2) DEFAULT 0,
    price_range_max NUMERIC(12,2) DEFAULT 0,
    preferred_price_mid NUMERIC(12,2) DEFAULT 0,
    total_interactions INT DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_similarity (
    product_id VARCHAR(255) NOT NULL,
    similar_product_id VARCHAR(255) NOT NULL,
    similarity_score NUMERIC(5,4) NOT NULL,
    method VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (product_id, similar_product_id, method)
);

CREATE TABLE IF NOT EXISTS trending_scores (
    product_id VARCHAR(255) PRIMARY KEY,
    score NUMERIC(5,4) NOT NULL DEFAULT 0,
    velocity NUMERIC(10,4) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rec_events_user ON recommendation_events(user_id);
CREATE INDEX IF NOT EXISTS idx_rec_events_product ON recommendation_events(product_id);
CREATE INDEX IF NOT EXISTS idx_rec_events_type ON recommendation_events(event_type);
CREATE INDEX IF NOT EXISTS idx_rec_events_created ON recommendation_events(created_at);
CREATE INDEX IF NOT EXISTS idx_product_similarity_score ON product_similarity(similarity_score DESC);
CREATE INDEX IF NOT EXISTS idx_trending_score ON trending_scores(score DESC);
