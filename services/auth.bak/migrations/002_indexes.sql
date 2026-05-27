-- Migration: 002_indexes
-- Description: Additional indexes for performance

ALTER TABLE audit_logs ADD INDEX idx_audit_ip (ip_address);
ALTER TABLE sessions ADD INDEX idx_sessions_device (device_id);
