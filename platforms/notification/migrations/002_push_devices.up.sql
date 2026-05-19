CREATE TABLE IF NOT EXISTS push_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    token TEXT NOT NULL,
    platform VARCHAR(50) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_push_devices_user_id ON push_devices(user_id);
CREATE INDEX idx_push_devices_token ON push_devices(token);
CREATE UNIQUE INDEX idx_push_devices_user_token ON push_devices(user_id, token);
