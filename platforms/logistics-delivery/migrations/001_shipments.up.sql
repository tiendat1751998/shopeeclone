CREATE TABLE IF NOT EXISTS shipments (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL,
    customer_id TEXT NOT NULL,
    warehouse_id TEXT NOT NULL DEFAULT '',
    courier_id TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    origin_street TEXT NOT NULL DEFAULT '',
    origin_city TEXT NOT NULL DEFAULT '',
    origin_state TEXT NOT NULL DEFAULT '',
    origin_country TEXT NOT NULL DEFAULT '',
    origin_zip TEXT NOT NULL DEFAULT '',
    origin_lat DOUBLE PRECISION DEFAULT 0,
    origin_lng DOUBLE PRECISION DEFAULT 0,
    dest_street TEXT NOT NULL DEFAULT '',
    dest_city TEXT NOT NULL DEFAULT '',
    dest_state TEXT NOT NULL DEFAULT '',
    dest_country TEXT NOT NULL DEFAULT '',
    dest_zip TEXT NOT NULL DEFAULT '',
    dest_lat DOUBLE PRECISION DEFAULT 0,
    dest_lng DOUBLE PRECISION DEFAULT 0,
    total_weight DOUBLE PRECISION DEFAULT 0,
    total_volume DOUBLE PRECISION DEFAULT 0,
    estimated_distance DOUBLE PRECISION DEFAULT 0,
    estimated_eta TIMESTAMPTZ,
    actual_delivered_at TIMESTAMPTZ,
    courier_notes TEXT DEFAULT '',
    replay_id TEXT DEFAULT '',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_customer_id ON shipments(customer_id);
CREATE INDEX idx_shipments_courier_id ON shipments(courier_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipments_warehouse_id ON shipments(warehouse_id);
CREATE INDEX idx_shipments_replay_id ON shipments(replay_id);
