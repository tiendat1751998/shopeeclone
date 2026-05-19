-- Notification Campaign Service - Initial Schema
-- Campaign management, templates, and delivery tracking

CREATE TABLE IF NOT EXISTS notification_campaigns (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    campaign_type ENUM('promotional','transactional','system','retargeting','abandoned_cart','price_drop','back_in_stock') NOT NULL DEFAULT 'promotional',
    channel ENUM('email','sms','push','in_app','whatsapp','all') NOT NULL DEFAULT 'email',
    status ENUM('draft','scheduled','active','paused','completed','cancelled') NOT NULL DEFAULT 'draft',
    priority ENUM('low','normal','high','urgent') NOT NULL DEFAULT 'normal',
    target_audience JSON DEFAULT NULL,
    segment_criteria JSON DEFAULT NULL,
    scheduled_at TIMESTAMP NULL DEFAULT NULL,
    started_at TIMESTAMP NULL DEFAULT NULL,
    completed_at TIMESTAMP NULL DEFAULT NULL,
    total_recipients INT NOT NULL DEFAULT 0,
    sent_count INT NOT NULL DEFAULT 0,
    delivered_count INT NOT NULL DEFAULT 0,
    opened_count INT NOT NULL DEFAULT 0,
    clicked_count INT NOT NULL DEFAULT 0,
    bounced_count INT NOT NULL DEFAULT 0,
    unsubscribed_count INT NOT NULL DEFAULT 0,
    created_by VARCHAR(36) DEFAULT NULL,
    approved_by VARCHAR(36) DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_nc_status (status),
    INDEX idx_nc_type (campaign_type, status),
    INDEX idx_nc_scheduled (scheduled_at, status),
    INDEX idx_nc_channel (channel, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS notification_templates (
    id VARCHAR(36) PRIMARY KEY,
    campaign_id VARCHAR(36) DEFAULT NULL,
    name VARCHAR(255) NOT NULL,
    channel ENUM('email','sms','push','in_app','whatsapp') NOT NULL DEFAULT 'email',
    subject VARCHAR(500) DEFAULT NULL,
    body_html TEXT DEFAULT NULL,
    body_text TEXT DEFAULT NULL,
    body_push VARCHAR(500) DEFAULT NULL,
    variables_schema JSON DEFAULT NULL,
    locale VARCHAR(10) NOT NULL DEFAULT 'en',
    version INT NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_nt_campaign (campaign_id),
    INDEX idx_nt_channel (channel, locale, is_active),
    INDEX idx_nt_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS campaign_recipients (
    id VARCHAR(36) PRIMARY KEY,
    campaign_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    channel ENUM('email','sms','push','in_app','whatsapp') NOT NULL DEFAULT 'email',
    status ENUM('pending','sent','delivered','opened','clicked','bounced','failed','unsubscribed') NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMP NULL DEFAULT NULL,
    delivered_at TIMESTAMP NULL DEFAULT NULL,
    opened_at TIMESTAMP NULL DEFAULT NULL,
    clicked_at TIMESTAMP NULL DEFAULT NULL,
    error_message TEXT DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_cr_campaign_user (campaign_id, user_id, channel),
    INDEX idx_cr_status (status),
    INDEX idx_cr_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS notification_preferences (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    email_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sms_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    push_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    in_app_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    whatsapp_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    promotional_email BOOLEAN NOT NULL DEFAULT TRUE,
    transactional_email BOOLEAN NOT NULL DEFAULT TRUE,
    order_updates_sms BOOLEAN NOT NULL DEFAULT TRUE,
    price_alerts_push BOOLEAN NOT NULL DEFAULT TRUE,
    quiet_hours_start TIME DEFAULT NULL,
    quiet_hours_end TIME DEFAULT NULL,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_np_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS notification_delivery_logs (
    id VARCHAR(36) PRIMARY KEY,
    campaign_id VARCHAR(36) DEFAULT NULL,
    template_id VARCHAR(36) DEFAULT NULL,
    user_id VARCHAR(36) NOT NULL,
    channel ENUM('email','sms','push','in_app','whatsapp') NOT NULL,
    status ENUM('queued','sent','delivered','opened','clicked','bounced','failed','suppressed') NOT NULL DEFAULT 'queued',
    provider VARCHAR(64) DEFAULT NULL,
    provider_message_id VARCHAR(255) DEFAULT NULL,
    error_code VARCHAR(64) DEFAULT NULL,
    error_message TEXT DEFAULT NULL,
    retry_count INT NOT NULL DEFAULT 0,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ndl_campaign (campaign_id),
    INDEX idx_ndl_user (user_id),
    INDEX idx_ndl_status (status),
    INDEX idx_ndl_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
