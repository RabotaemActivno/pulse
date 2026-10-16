BEGIN;

CREATE TABLE IF NOT EXISTS checks
(
    id UUID PRIMARY KEY,
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL,
    http_status INT,
    latency_ms INT NOT NULL,
    error TEXT
);

COMMIT;