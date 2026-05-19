CREATE TABLE IF NOT EXISTS fulfillments (
    id TEXT PRIMARY KEY,
    shipment_id TEXT NOT NULL,
    order_id TEXT NOT NULL,
    warehouse_id TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    packed_at TIMESTAMPTZ,
    shipped_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    notes TEXT DEFAULT '',
    replay_id TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fulfillments_shipment ON fulfillments(shipment_id);
CREATE INDEX idx_fulfillments_order ON fulfillments(order_id);
CREATE INDEX idx_fulfillments_status ON fulfillments(status);

CREATE TABLE IF NOT EXISTS pickup_schedule (
    id TEXT PRIMARY KEY,
    shipment_id TEXT NOT NULL,
    fulfillment_id TEXT NOT NULL DEFAULT '',
    courier_id TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'scheduled',
    address TEXT DEFAULT '',
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,
    scheduled_at TIMESTAMPTZ,
    picked_up_at TIMESTAMPTZ,
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pickup_shipment ON pickup_schedule(shipment_id);
CREATE INDEX idx_pickup_status ON pickup_schedule(status);
