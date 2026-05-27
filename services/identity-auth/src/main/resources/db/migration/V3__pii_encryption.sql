-- PII Encryption Support
-- Migration: V3
-- Description: Add email_hash column for encrypted email lookups, update column sizes for encrypted data

ALTER TABLE users
    ADD COLUMN email_hash VARCHAR(64) NOT NULL DEFAULT '' AFTER email,
    MODIFY COLUMN email VARCHAR(512) NOT NULL,
    MODIFY COLUMN phone VARCHAR(256),
    MODIFY COLUMN full_name VARCHAR(512) NOT NULL;

CREATE INDEX idx_users_email_hash ON users(email_hash);

-- Migrate existing data: populate email_hash for existing records
UPDATE users SET email_hash = LOWER(CONCAT(
    HEX(UNHEX(SHA2(email, 256)))
)) WHERE email_hash = '';
