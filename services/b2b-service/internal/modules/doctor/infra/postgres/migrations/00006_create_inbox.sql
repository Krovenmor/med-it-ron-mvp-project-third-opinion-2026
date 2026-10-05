-- +goose Up
CREATE TABLE inbox (
    message_id   text        PRIMARY KEY,
    topic        text        NOT NULL,
    processed_at timestamptz NOT NULL
);

-- +goose Down
DROP TABLE inbox;
