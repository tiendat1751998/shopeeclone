CREATE TABLE IF NOT EXISTS campaigns (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    campaign_type VARCHAR(10) NOT NULL,
    daily_budget NUMERIC(12,2) DEFAULT 0,
    lifetime_budget NUMERIC(12,2) DEFAULT 0,
    bid_amount NUMERIC(10,4) DEFAULT 0,
    target_cpa NUMERIC(10,4) DEFAULT 0,
    quality_score NUMERIC(3,2) DEFAULT 1.0,
    date_start TIMESTAMPTZ,
    date_end TIMESTAMPTZ,
    targeting JSONB DEFAULT '{}',
    creative_ids JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS creatives (
    id VARCHAR(36) PRIMARY KEY,
    campaign_id VARCHAR(36) NOT NULL REFERENCES campaigns(id),
    name VARCHAR(255) NOT NULL,
    format VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    content TEXT,
    destination_url TEXT,
    sizes JSONB DEFAULT '[]',
    impressions BIGINT DEFAULT 0,
    clicks BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS impressions (
    id VARCHAR(36) PRIMARY KEY,
    campaign_id VARCHAR(36) NOT NULL,
    creative_id VARCHAR(36),
    user_id VARCHAR(255),
    cost NUMERIC(10,4) DEFAULT 0,
    device VARCHAR(50),
    location VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS clicks (
    id VARCHAR(36) PRIMARY KEY,
    impression_id VARCHAR(36) REFERENCES impressions(id),
    campaign_id VARCHAR(36) NOT NULL,
    creative_id VARCHAR(36),
    user_id VARCHAR(255),
    cost NUMERIC(10,4) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS conversions (
    id VARCHAR(36) PRIMARY KEY,
    click_id VARCHAR(36) REFERENCES clicks(id),
    campaign_id VARCHAR(36) NOT NULL,
    creative_id VARCHAR(36),
    user_id VARCHAR(255),
    revenue NUMERIC(12,2) DEFAULT 0,
    conversion_type VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bid_history (
    id BIGSERIAL PRIMARY KEY,
    campaign_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(255),
    bid_amount NUMERIC(10,4) NOT NULL,
    won BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS budget_trackers (
    campaign_id VARCHAR(36) PRIMARY KEY,
    daily_budget NUMERIC(12,2) DEFAULT 0,
    lifetime_budget NUMERIC(12,2) DEFAULT 0,
    spent_today NUMERIC(12,2) DEFAULT 0,
    spent_total NUMERIC(12,2) DEFAULT 0,
    last_reset_date DATE,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_campaigns_status ON campaigns(status);
CREATE INDEX IF NOT EXISTS idx_creatives_campaign ON creatives(campaign_id);
CREATE INDEX IF NOT EXISTS idx_impressions_campaign ON impressions(campaign_id);
CREATE INDEX IF NOT EXISTS idx_clicks_campaign ON clicks(campaign_id);
CREATE INDEX IF NOT EXISTS idx_conversions_campaign ON conversions(campaign_id);
CREATE INDEX IF NOT EXISTS idx_bid_history_campaign ON bid_history(campaign_id);
CREATE INDEX IF NOT EXISTS idx_impressions_created ON impressions(created_at);
CREATE INDEX IF NOT EXISTS idx_clicks_created ON clicks(created_at);
CREATE INDEX IF NOT EXISTS idx_conversions_created ON conversions(created_at);
