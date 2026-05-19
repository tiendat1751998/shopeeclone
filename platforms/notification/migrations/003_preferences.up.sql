CREATE TABLE IF NOT EXISTS user_preferences (
    user_id VARCHAR(255) PRIMARY KEY,
    channel_push BOOLEAN NOT NULL DEFAULT true,
    channel_email BOOLEAN NOT NULL DEFAULT true,
    channel_sms BOOLEAN NOT NULL DEFAULT true,
    channel_inapp BOOLEAN NOT NULL DEFAULT true,
    categories JSONB DEFAULT '{}',
    quiet_hours_enabled BOOLEAN NOT NULL DEFAULT false,
    quiet_hours_start VARCHAR(5) DEFAULT '22:00',
    quiet_hours_end VARCHAR(5) DEFAULT '08:00',
    quiet_hours_timezone VARCHAR(50) DEFAULT 'UTC',
    email_digest BOOLEAN NOT NULL DEFAULT false,
    push_enabled BOOLEAN NOT NULL DEFAULT true,
    sms_promotions BOOLEAN NOT NULL DEFAULT true,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS suppression_list (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    reason VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_suppression_user_id ON suppression_list(user_id);
CREATE INDEX idx_suppression_email ON suppression_list(email);
CREATE INDEX idx_suppression_phone ON suppression_list(phone);
