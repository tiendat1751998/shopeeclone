CREATE TABLE IF NOT EXISTS tracking_events (
    id TEXT PRIMARY KEY,
    shipment_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    lat DOUBLE PRECISION DEFAULT 0,
    lng DOUBLE PRECISION DEFAULT 0,
    location_name TEXT DEFAULT '',
    location_address TEXT DEFAULT '',
    description TEXT DEFAULT '',
    courier_data JSONB DEFAULT '{}',
    replay_id TEXT DEFAULT '',
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tracking_events_shipment ON tracking_events(shipment_id);
CREATE INDEX idx_tracking_events_type ON tracking_events(event_type);
CREATE INDEX idx_tracking_events_occurred ON tracking_events(occurred_at DESC);
CREATE INDEX idx_tracking_events_replay ON tracking_events(replay_id);
