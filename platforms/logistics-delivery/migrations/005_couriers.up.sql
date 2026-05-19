CREATE TABLE IF NOT EXISTS couriers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL DEFAULT '',
    provider TEXT NOT NULL DEFAULT 'internal',
    status TEXT NOT NULL DEFAULT 'available',
    zone_id TEXT DEFAULT '',
    current_lat DOUBLE PRECISION DEFAULT 0,
    current_lng DOUBLE PRECISION DEFAULT 0,
    last_seen_at TIMESTAMPTZ,
    max_capacity INT DEFAULT 10,
    current_load INT DEFAULT 0,
    rating DOUBLE PRECISION DEFAULT 5.0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_couriers_zone ON couriers(zone_id);
CREATE INDEX idx_couriers_status ON couriers(status);
CREATE INDEX idx_couriers_active ON couriers(is_active);
