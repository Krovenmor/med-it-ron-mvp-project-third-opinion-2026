-- +goose Up
CREATE TABLE case_events (
    id          bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    case_id     uuid        NOT NULL REFERENCES routes (case_id),
    type        text        NOT NULL,
    from_status text,
    to_status   text,
    actor       text        NOT NULL,
    payload     jsonb       NOT NULL DEFAULT '{}',
    occurred_at timestamptz NOT NULL
);

CREATE INDEX case_events_case_id_idx ON case_events (case_id, occurred_at);

-- +goose Down
DROP TABLE case_events;
