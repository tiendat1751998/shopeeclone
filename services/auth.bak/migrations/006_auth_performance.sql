-- ============================================================
-- Migration 006: Auth Performance Optimization
-- Fixes the #1 bottleneck: encrypted email lookup does full table scan
-- ============================================================

-- PROBLEM: The users.email column stores AES-encrypted values.
-- When auth service does WHERE email = ?, MySQL must scan every row,
-- decrypt each email, and compare. With 10K users = 10K decryptions per login.
--
-- FIX: Add email_hash column (SHA-256 of plaintext email) and index it.
-- Login query becomes: WHERE email_hash = SHA2(?, 256) AND email = ?
-- The hash lookup uses the index (O(log n)), then only 1 row is decrypted.

-- Add email_hash column if not exists
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_hash VARCHAR(64) DEFAULT NULL AFTER email;

-- Add index on email_hash for O(log n) login lookups
ALTER TABLE users ADD INDEX idx_users_email_hash (email_hash);

-- Create a covering index for login queries (avoids table lookup)
-- Only fetches the columns needed for login: password_hash, status, id
ALTER TABLE users ADD INDEX idx_users_login_cover (email_hash, email, password_hash, status, id, display_name, failed_attempts, locked_until);

-- Update existing rows: populate email_hash from plaintext email
-- NOTE: This only works for users where email is still decryptable.
-- For the seed data, we populate it via the application on next login.
-- For new registrations, the hash is set automatically.

-- Add index on sessions for faster active session lookups during login
ALTER TABLE sessions ADD INDEX idx_sessions_user_active (user_id, status, expires_at);

-- ============================================================
-- Partial index approach for active sessions (MySQL 8.0+)
-- This makes FindActiveByUserID() much faster
-- ============================================================

-- Composite index: WHERE user_id=? AND status='active' AND expires_at > NOW()
ALTER TABLE sessions ADD INDEX idx_sessions_active_lookup (user_id, status, last_active_at DESC);
