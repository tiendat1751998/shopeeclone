-- Migration: 002_indexes
-- Description: Additional indexes for performance

ALTER TABLE audit_logs ADD INDEX idx_audit_action_created (action, created_at);
ALTER TABLE sessions ADD INDEX idx_sessions_user_status (user_id, status);
ALTER TABLE login_attempts ADD INDEX idx_login_attempts_email_created (email, created_at);
ALTER TABLE audit_logs ADD INDEX idx_audit_ip (ip);
ALTER TABLE sessions ADD INDEX idx_sessions_device (device_id);
