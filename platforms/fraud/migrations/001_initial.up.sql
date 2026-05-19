-- Fraud Detection Platform - PostgreSQL Schema
-- Migration 001: Initial schema

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS fraud_rules (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    condition JSONB NOT NULL,
    severity INT NOT NULL DEFAULT 1,
    weight FLOAT NOT NULL DEFAULT 1.0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    cooldown INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fraud_rules_active ON fraud_rules(is_active);

CREATE TABLE IF NOT EXISTS fraud_alerts (
    id VARCHAR(64) PRIMARY KEY,
    event_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    type VARCHAR(32) NOT NULL,
    risk_score FLOAT NOT NULL,
    risk_level VARCHAR(16) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    resolved_by VARCHAR(64),
    resolution TEXT
);

CREATE INDEX idx_fraud_alerts_user ON fraud_alerts(user_id);
CREATE INDEX idx_fraud_alerts_status ON fraud_alerts(status);
CREATE INDEX idx_fraud_alerts_type ON fraud_alerts(type);
CREATE INDEX idx_fraud_alerts_created ON fraud_alerts(created_at DESC);

CREATE TABLE IF NOT EXISTS fraud_cases (
    id VARCHAR(64) PRIMARY KEY,
    alert_id VARCHAR(64) NOT NULL REFERENCES fraud_alerts(id),
    user_id VARCHAR(64) NOT NULL,
    title VARCHAR(256) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    priority VARCHAR(16) NOT NULL DEFAULT 'medium',
    risk_score FLOAT NOT NULL DEFAULT 0,
    investigator VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    resolution TEXT
);

CREATE INDEX idx_fraud_cases_user ON fraud_cases(user_id);
CREATE INDEX idx_fraud_cases_status ON fraud_cases(status);
CREATE INDEX idx_fraud_cases_priority ON fraud_cases(priority);

CREATE TABLE IF NOT EXISTS fraud_evidence (
    id VARCHAR(64) PRIMARY KEY,
    case_id VARCHAR(64) NOT NULL REFERENCES fraud_cases(id),
    type VARCHAR(32) NOT NULL,
    description TEXT,
    data TEXT,
    added_by VARCHAR(64),
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fraud_evidence_case ON fraud_evidence(case_id);

CREATE TABLE IF NOT EXISTS blacklist_entries (
    id VARCHAR(64) PRIMARY KEY,
    type VARCHAR(16) NOT NULL,
    value VARCHAR(256) NOT NULL,
    reason VARCHAR(64),
    added_by VARCHAR(64),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_blacklist_type_value ON blacklist_entries(type, value);
CREATE INDEX idx_blacklist_active ON blacklist_entries(is_active);

CREATE TABLE IF NOT EXISTS verification_requests (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    method VARCHAR(16) NOT NULL,
    target VARCHAR(256) NOT NULL,
    code_hash VARCHAR(128) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ
);

CREATE INDEX idx_verification_user ON verification_requests(user_id, method);
CREATE INDEX idx_verification_status ON verification_requests(status);

CREATE TABLE IF NOT EXISTS kyc_statuses (
    user_id VARCHAR(64) PRIMARY KEY,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    level VARCHAR(16) NOT NULL DEFAULT 'basic',
    document_type VARCHAR(32),
    submitted_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    rejected_reason TEXT
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id VARCHAR(64) NOT NULL,
    action VARCHAR(64) NOT NULL,
    resource VARCHAR(64) NOT NULL,
    resource_id VARCHAR(64) NOT NULL,
    details JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_actor ON audit_logs(actor_id, created_at DESC);
CREATE INDEX idx_audit_resource ON audit_logs(resource, resource_id, created_at DESC);
