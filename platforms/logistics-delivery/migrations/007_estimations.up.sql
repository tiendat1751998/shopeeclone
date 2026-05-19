CREATE TABLE IF NOT EXISTS estimations (
    id TEXT PRIMARY KEY,
    shipment_id TEXT NOT NULL,
    distance_km DOUBLE PRECISION DEFAULT 0,
    base_duration_min INT DEFAULT 0,
    traffic_delay_min INT DEFAULT 0,
    weather_delay_min INT DEFAULT 0,
    total_duration_min INT DEFAULT 0,
    eta TIMESTAMPTZ NOT NULL,
    confidence DOUBLE PRECISION DEFAULT 0.5,
    route_hash TEXT DEFAULT '',
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_estimations_shipment ON estimations(shipment_id);
CREATE INDEX idx_estimations_hash ON estimations(route_hash);
