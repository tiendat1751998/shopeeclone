CREATE TABLE IF NOT EXISTS notifications (
    id VARCHAR(36) PRIMARY KEY, user_id VARCHAR(36) NOT NULL, type VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL, body TEXT, data TEXT, channel VARCHAR(20) NOT NULL,
    status ENUM('pending','sent','delivered','failed','read') NOT NULL DEFAULT 'pending',
    priority INT NOT NULL DEFAULT 0, read_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_notifications_user (user_id, status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS notification_templates (
    id VARCHAR(36) PRIMARY KEY, name VARCHAR(100) NOT NULL UNIQUE, type VARCHAR(20) NOT NULL,
    subject VARCHAR(500) NOT NULL, body TEXT NOT NULL, variables TEXT,
    version INT NOT NULL DEFAULT 1, is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id VARCHAR(36) NOT NULL, channel VARCHAR(20) NOT NULL, enabled BOOLEAN DEFAULT TRUE,
    quiet_hours VARCHAR(50) DEFAULT NULL, PRIMARY KEY (user_id, channel)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS delivery_logs (
    id VARCHAR(36) PRIMARY KEY, notification_id VARCHAR(36) NOT NULL, channel VARCHAR(20) NOT NULL,
    provider VARCHAR(50) NOT NULL, status VARCHAR(20) NOT NULL, error_message TEXT,
    attempt_count INT NOT NULL DEFAULT 0, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_delivery_notif (notification_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
