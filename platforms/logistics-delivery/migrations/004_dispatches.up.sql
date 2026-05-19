CREATE TABLE IF NOT EXISTS dispatches (
    id TEXT PRIMARY KEY,
    shipment_id TEXT NOT NULL,
    courier_id TEXT DEFAULT '',
    zone_id TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    pickup_time TIMESTAMPTZ,
    dispatch_time TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dispatches_shipment ON dispatches(shipment_id);
CREATE INDEX idx_dispatches_courier ON dispatches(courier_id);
CREATE INDEX idx_dispatches_status ON dispatches(status);
