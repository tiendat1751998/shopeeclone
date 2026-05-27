-- ============================================================
-- Migration: Shipment Service Performance Indexes
-- Adds covering indexes for hot-path queries
-- ============================================================

-- Covering index for shipment lookup by order_id
ALTER TABLE shipments ADD INDEX idx_shipments_order_cover (order_id, deleted_at, id, status, carrier_id, tracking_number, cost, currency, created_at);

-- Covering index for shipment lookup by user
ALTER TABLE shipments ADD INDEX idx_shipments_user_cover (user_id, deleted_at, status, created_at DESC, id, order_id, carrier_id, tracking_number);

-- Covering index for tracking events
ALTER TABLE tracking_events ADD INDEX idx_tracking_shipment_cover (shipment_id, timestamp ASC, id, status, location, description);

-- Covering index for scan events by shipment
ALTER TABLE scan_events ADD INDEX idx_scan_shipment_cover (shipment_id, created_at DESC, id, qr_code_id, shipper_id, scan_type, is_valid);

-- Covering index for scan events by shipper
ALTER TABLE scan_events ADD INDEX idx_scan_shipper_cover (shipper_id, created_at DESC, id, shipment_id, qr_code_id, scan_type);

-- Covering index for qr_codes lookup
ALTER TABLE qr_codes ADD INDEX idx_qr_codes_code_cover (code, id, shipment_id, type, status, expires_at);

-- Covering index for idempotency keys
ALTER TABLE idempotency_keys ADD INDEX idx_idempotency_key_cover (`key`, shipment_id, expires_at);

ANALYZE TABLE shipments;
ANALYZE TABLE tracking_events;
ANALYZE TABLE scan_events;
ANALYZE TABLE qr_codes;
ANALYZE TABLE idempotency_keys;
