-- Analytics Platform - PostgreSQL Schema
-- Migration 001: Initial schema

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id VARCHAR(64) NOT NULL UNIQUE,
    event_type VARCHAR(32) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    session_id VARCHAR(64),
    timestamp TIMESTAMPTZ NOT NULL,
    properties JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    device VARCHAR(32),
    referrer TEXT,
    source VARCHAR(32),
    country VARCHAR(8),
    campaign VARCHAR(64),
    revenue DECIMAL(20,4) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analytics_events_event_type ON analytics_events(event_type);
CREATE INDEX idx_analytics_events_user_id ON analytics_events(user_id);
CREATE INDEX idx_analytics_events_timestamp ON analytics_events(timestamp DESC);
CREATE INDEX idx_analytics_events_session ON analytics_events(session_id);
CREATE INDEX idx_analytics_events_source ON analytics_events(source);
CREATE INDEX idx_analytics_events_country ON analytics_events(country);

CREATE TABLE IF NOT EXISTS analytics_sessions (
    session_id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    duration_seconds BIGINT DEFAULT 0,
    pageviews INT DEFAULT 0,
    events_count INT DEFAULT 0,
    device VARCHAR(32),
    source VARCHAR(32),
    country VARCHAR(8),
    is_active BOOLEAN DEFAULT true,
    has_conversion BOOLEAN DEFAULT false,
    revenue DECIMAL(20,4) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analytics_sessions_user ON analytics_sessions(user_id);
CREATE INDEX idx_analytics_sessions_start ON analytics_sessions(start_time DESC);
CREATE INDEX idx_analytics_sessions_active ON analytics_sessions(is_active) WHERE is_active = true;

CREATE TABLE IF NOT EXISTS analytics_reports (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    query JSONB NOT NULL,
    result JSONB,
    organization_id VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics_funnels (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    definition JSONB NOT NULL,
    result JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics_cohorts (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    definition JSONB NOT NULL,
    result JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics_dashboards (
    id VARCHAR(64) PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    description TEXT,
    organization_id VARCHAR(64) NOT NULL,
    created_by VARCHAR(64) NOT NULL,
    is_public BOOLEAN DEFAULT false,
    tags TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics_widgets (
    id VARCHAR(64) PRIMARY KEY,
    dashboard_id VARCHAR(64) NOT NULL REFERENCES analytics_dashboards(id) ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL,
    type VARCHAR(32) NOT NULL,
    width INT DEFAULT 4,
    height INT DEFAULT 3,
    position_x INT DEFAULT 0,
    position_y INT DEFAULT 0,
    data_source JSONB NOT NULL,
    config JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analytics_widgets_dashboard ON analytics_widgets(dashboard_id);

CREATE TABLE IF NOT EXISTS analytics_scheduled_reports (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    query JSONB NOT NULL,
    frequency VARCHAR(16) NOT NULL,
    delivery_channel VARCHAR(16) NOT NULL,
    recipients TEXT[],
    webhook_url TEXT,
    format VARCHAR(16) DEFAULT 'csv',
    time_zone VARCHAR(32) DEFAULT 'UTC',
    next_run_at TIMESTAMPTZ NOT NULL,
    last_run_at TIMESTAMPTZ,
    last_status VARCHAR(32),
    is_active BOOLEAN DEFAULT true,
    created_by VARCHAR(64),
    organization_id VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analytics_scheduled_next_run ON analytics_scheduled_reports(next_run_at) WHERE is_active = true;

CREATE TABLE IF NOT EXISTS analytics_report_generations (
    id VARCHAR(64) PRIMARY KEY,
    report_id VARCHAR(64) NOT NULL REFERENCES analytics_scheduled_reports(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    data JSONB,
    error TEXT,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delivered_at TIMESTAMPTZ
);

CREATE INDEX idx_analytics_generations_report ON analytics_report_generations(report_id);
