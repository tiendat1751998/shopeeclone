-- ============================================================
-- Migration: 002_outbox_three_state.sql
-- Description: Add three-state outbox columns to payment service
-- ============================================================

ALTER TABLE outbox_events
  ADD COLUMN error_message TEXT AFTER payload,
  ADD COLUMN processing_at TIMESTAMP NULL DEFAULT NULL AFTER error_message,
  ADD INDEX idx_outbox_created (created_at);
