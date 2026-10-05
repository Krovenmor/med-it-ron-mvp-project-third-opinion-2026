-- +goose Up
CREATE TABLE jobs (
    id           bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    kind         text        NOT NULL,
    key          text        NOT NULL,
    payload      jsonb       NOT NULL DEFAULT '{}',
    state        text        NOT NULL DEFAULT 'pending' CHECK (state IN ('pending', 'done', 'failed', 'cancelled')),
    run_at       timestamptz NOT NULL,
    locked_until timestamptz,
    attempts     int         NOT NULL DEFAULT 0,
    last_error   text        NOT NULL DEFAULT '',
    created_at   timestamptz NOT NULL
);

CREATE INDEX jobs_pending_run_at_idx ON jobs (run_at) WHERE state = 'pending';
CREATE INDEX jobs_key_idx ON jobs (key);

-- +goose Down
DROP TABLE jobs;
