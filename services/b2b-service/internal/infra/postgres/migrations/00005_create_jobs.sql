-- +goose Up
CREATE TABLE jobs (
    id           bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    case_id      uuid        NOT NULL REFERENCES cases (id),
    kind         text        NOT NULL,
    state        text        NOT NULL DEFAULT 'pending' CHECK (state IN ('pending', 'done', 'failed', 'cancelled')),
    run_at       timestamptz NOT NULL,
    locked_until timestamptz,
    attempts     int         NOT NULL DEFAULT 0,
    last_error   text        NOT NULL DEFAULT '',
    created_at   timestamptz NOT NULL
);

CREATE INDEX jobs_pending_run_at_idx ON jobs (run_at) WHERE state = 'pending';

-- +goose Down
DROP TABLE jobs;
