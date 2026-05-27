-- ============================================================
-- Migration: Notification Service Performance Indexes
-- Adds covering indexes for hot-path queries
-- ============================================================

-- Covering index for notification listing by user (avoids table access)
ALTER TABLE notifications ADD INDEX idx_notifs_user_cover (user_id, created_at DESC, id, type, title, body, data, channel, status, priority);

-- Covering index for unread notifications
ALTER TABLE notifications ADD INDEX idx_notifs_unread_cover (user_id, status, created_at DESC, id, type, title);

ANALYZE TABLE notifications;
