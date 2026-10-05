-- +goose Up
CREATE TABLE patients (
    id            uuid        PRIMARY KEY,
    source_system text        NOT NULL,
    external_id   text        NOT NULL,
    full_name     text        NOT NULL,
    birth_date    date        NOT NULL,
    sex           text        NOT NULL,
    phone         text        NOT NULL DEFAULT '',
    email         text        NOT NULL DEFAULT '',
    updated_at    timestamptz NOT NULL,
    UNIQUE (source_system, external_id)
);

-- +goose Down
DROP TABLE patients;
