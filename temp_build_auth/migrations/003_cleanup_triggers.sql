-- Migration: 003_cleanup_triggers
-- Description: Automated cleanup triggers

CREATE EVENT IF NOT EXISTS cleanup_expired_sessions
ON SCHEDULE EVERY 1 HOUR
DO
    DELETE FROM sessions WHERE expires_at < NOW() OR status IN ('expired', 'revoked');

CREATE EVENT IF NOT EXISTS cleanup_old_audit_logs
ON SCHEDULE EVERY 1 DAY
DO
    DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL 90 DAY;

CREATE EVENT IF NOT EXISTS cleanup_old_login_attempts
ON SCHEDULE EVERY 1 DAY
DO
    DELETE FROM login_attempts WHERE created_at < NOW() - INTERVAL 30 DAY;

CREATE EVENT IF NOT EXISTS cleanup_revoked_refresh_tokens
ON SCHEDULE EVERY 1 HOUR
DO
    DELETE FROM refresh_tokens WHERE expires_at < NOW() OR revoked = TRUE;
