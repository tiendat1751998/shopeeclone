-- Live Commerce Platform - ClickHouse Schema
-- Migration 001: Analytics tables

CREATE TABLE IF NOT EXISTS live_viewer_events (
    room_id String,
    user_id String,
    event_type String,
    timestamp DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (room_id, timestamp)
TTL timestamp + INTERVAL 90 DAY;

CREATE TABLE IF NOT EXISTS live_engagement_events (
    room_id String,
    user_id String,
    event_type String,
    value Int64,
    timestamp DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (room_id, event_type, timestamp)
TTL timestamp + INTERVAL 90 DAY;

CREATE MATERIALIZED VIEW IF NOT EXISTS live_daily_stats
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMMDD(day)
ORDER BY (room_id, day)
AS SELECT
    room_id,
    toDate(timestamp) AS day,
    countDistinctIf(user_id, event_type = 'join') AS unique_viewers,
    countIf(event_type = 'chat') AS chat_messages,
    countIf(event_type = 'gift') AS gift_count,
    sumIf(value, event_type = 'gift') AS gift_value,
    countIf(event_type = 'reaction_like') AS likes,
    countIf(event_type = 'reaction_love') AS loves
FROM live_engagement_events
GROUP BY room_id, day;

CREATE TABLE IF NOT EXISTS live_concurrent_viewers (
    room_id String,
    viewer_count Int32,
    timestamp DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (room_id, timestamp)
TTL timestamp + INTERVAL 7 DAY;
