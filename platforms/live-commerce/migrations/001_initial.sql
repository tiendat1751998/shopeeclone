-- Live Commerce Platform - PostgreSQL Schema
-- Migration 001: Initial schema

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS livestreams (
    id VARCHAR(64) PRIMARY KEY,
    seller_id VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    cover_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    category VARCHAR(64),
    tags TEXT[] DEFAULT '{}',
    viewer_count BIGINT NOT NULL DEFAULT 0,
    peak_viewers BIGINT NOT NULL DEFAULT 0,
    total_likes BIGINT NOT NULL DEFAULT 0,
    total_gifts BIGINT NOT NULL DEFAULT 0,
    total_shares BIGINT NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    scheduled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_livestreams_seller ON livestreams(seller_id);
CREATE INDEX idx_livestreams_status ON livestreams(status);
CREATE INDEX idx_livestreams_status_viewers ON livestreams(status, viewer_count DESC) WHERE status = 'live';
CREATE INDEX idx_livestreams_category ON livestreams(category);
CREATE INDEX idx_livestreams_scheduled ON livestreams(scheduled_at) WHERE scheduled_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS chat_messages (
    id VARCHAR(64) PRIMARY KEY,
    room_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    username VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'text',
    is_moderated BOOLEAN NOT NULL DEFAULT FALSE,
    sequence BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chat_room_seq ON chat_messages(room_id, sequence DESC);
CREATE INDEX idx_chat_room_ts ON chat_messages(room_id, timestamp DESC);
CREATE INDEX idx_chat_user ON chat_messages(user_id);

CREATE TABLE IF NOT EXISTS reactions (
    id VARCHAR(64) PRIMARY KEY,
    room_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    type VARCHAR(20) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reactions_room ON reactions(room_id, type);
CREATE INDEX idx_reactions_room_ts ON reactions(room_id, timestamp DESC);

CREATE TABLE IF NOT EXISTS gifts (
    id VARCHAR(64) PRIMARY KEY,
    room_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    username VARCHAR(128) NOT NULL,
    gift_type VARCHAR(32) NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(8) NOT NULL DEFAULT 'VND',
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gifts_room ON gifts(room_id, timestamp DESC);
CREATE INDEX idx_gifts_room_amount ON gifts(room_id, amount DESC);
CREATE INDEX idx_gifts_user ON gifts(user_id);

CREATE TABLE IF NOT EXISTS pinned_products (
    id VARCHAR(64) PRIMARY KEY,
    livestream_id VARCHAR(64) NOT NULL REFERENCES livestreams(id) ON DELETE CASCADE,
    product_id VARCHAR(64) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL DEFAULT 0,
    image_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    pinned_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pinned_livestream ON pinned_products(livestream_id, is_active);

CREATE TABLE IF NOT EXISTS moderation_actions (
    id VARCHAR(64) PRIMARY KEY,
    room_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    action VARCHAR(32) NOT NULL,
    reason TEXT,
    moderated_by VARCHAR(64) NOT NULL,
    duration_sec BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mod_room ON moderation_actions(room_id, created_at DESC);
CREATE INDEX idx_mod_user ON moderation_actions(user_id);
CREATE INDEX idx_mod_action ON moderation_actions(action);

CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(64) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_status ON outbox_events(status, created_at);

CREATE TABLE IF NOT EXISTS viewer_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    username VARCHAR(128),
    connected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    disconnected_at TIMESTAMPTZ,
    duration_sec INT DEFAULT 0
);

CREATE INDEX idx_viewer_room ON viewer_sessions(room_id, connected_at DESC);
CREATE INDEX idx_viewer_user ON viewer_sessions(user_id);
