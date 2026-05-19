CREATE TABLE IF NOT EXISTS routes (
    id TEXT PRIMARY KEY,
    shipment_id TEXT NOT NULL,
    route_type TEXT NOT NULL DEFAULT 'warehouse',
    origin_id TEXT NOT NULL DEFAULT '',
    destination_id TEXT NOT NULL DEFAULT '',
    distance_km DOUBLE PRECISION DEFAULT 0,
    estimated_duration_min INT DEFAULT 0,
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_routes_shipment ON routes(shipment_id);
CREATE INDEX idx_routes_active ON routes(is_active);

CREATE TABLE IF NOT EXISTS zones (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_zones_city ON zones(city);
