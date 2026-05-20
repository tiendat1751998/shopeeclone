-- Order Service - Three-State Transactional Outbox
-- Migration: 003_outbox_three_state
-- Adds status tracking to prevent duplicate Kafka message delivery

ALTER TABLE outbox_events
  ADD COLUMN status ENUM('pending','processing','processed','failed') NOT NULL DEFAULT 'pending' AFTER event_type,
  ADD COLUMN error_message TEXT AFTER status,
  ADD COLUMN retries INT NOT NULL DEFAULT 0 AFTER error_message,
  ADD INDEX idx_outbox_status (status, created_at);

-- Migrate existing data: processed=TRUE → status='processed', processed=FALSE → status='pending'
UPDATE outbox_events SET status = 'processed' WHERE processed = TRUE;
UPDATE outbox_events SET status = 'pending' WHERE processed = FALSE;
