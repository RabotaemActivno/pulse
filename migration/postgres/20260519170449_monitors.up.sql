BEGIN;

CREATE TABLE IF NOT EXISTS monitors
(
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    url             TEXT NOT NULL,
    method          TEXT NOT NULL DEFAULT 'GET',
    interval_seconds INT NOT NULL,
    timeout_seconds  INT NOT NULL,
    expected_status  INT NOT NULL DEFAULT 200,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    next_check_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    consecutive_failures INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMIT;