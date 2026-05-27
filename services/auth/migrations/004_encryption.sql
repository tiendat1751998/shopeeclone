-- Auth Service - Encryption Support
-- Migration: 004_encryption
-- Description: Add hash columns for encrypted PII lookups

ALTER TABLE users
    ADD COLUMN email_hash VARCHAR(64) NOT NULL DEFAULT '' AFTER email,
    ADD COLUMN username_hash VARCHAR(64) NOT NULL DEFAULT '' AFTER username,
    ADD INDEX idx_users_email_hash (email_hash),
    ADD INDEX idx_users_username_hash (username_hash);
