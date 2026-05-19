CREATE TABLE IF NOT EXISTS sfu_nodes (
    id TEXT PRIMARY KEY,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    capacity INT NOT NULL DEFAULT 100,
    current_load INT NOT NULL DEFAULT 0,
    stream_count INT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active',
    last_heartbeat TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stream_sessions (
    id TEXT PRIMARY KEY,
    stream_id TEXT NOT NULL,
    node_id TEXT NOT NULL REFERENCES sfu_nodes(id),
    region TEXT NOT NULL,
    viewers INT NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cdn_endpoints (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    provider TEXT NOT NULL DEFAULT '',
    region TEXT NOT NULL,
    latency_ms INT NOT NULL DEFAULT 0,
    capacity INT NOT NULL DEFAULT 1000,
    current_load INT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cdn_purge_requests (
    id TEXT PRIMARY KEY,
    reason TEXT NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ws_cluster_nodes (
    id TEXT PRIMARY KEY,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    room_count INT NOT NULL DEFAULT 0,
    client_count INT NOT NULL DEFAULT 0,
    max_rooms INT NOT NULL DEFAULT 100,
    max_clients INT NOT NULL DEFAULT 1000,
    last_heartbeat TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS room_assignments (
    room_id TEXT PRIMARY KEY,
    node_id TEXT NOT NULL REFERENCES ws_cluster_nodes(id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stream_health (
    stream_id TEXT PRIMARY KEY,
    status TEXT NOT NULL DEFAULT 'healthy',
    node_id TEXT DEFAULT '',
    region TEXT DEFAULT '',
    last_checked TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS regions (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    latency_ms INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS region_latency (
    from_region TEXT NOT NULL,
    to_region TEXT NOT NULL,
    latency_ms INT NOT NULL DEFAULT 0,
    PRIMARY KEY (from_region, to_region)
);

CREATE TABLE IF NOT EXISTS transcode_jobs (
    id TEXT PRIMARY KEY,
    stream_id TEXT NOT NULL,
    input_url TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    error TEXT DEFAULT '',
    progress DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stream_sessions_node_id ON stream_sessions(node_id);
CREATE INDEX idx_stream_sessions_stream_id ON stream_sessions(stream_id);
CREATE INDEX idx_room_assignments_node_id ON room_assignments(node_id);
CREATE INDEX idx_transcode_jobs_stream_id ON transcode_jobs(stream_id);
