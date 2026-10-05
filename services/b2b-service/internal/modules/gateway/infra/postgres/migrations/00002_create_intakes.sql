-- +goose Up
CREATE TABLE intakes (
    case_id       uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id    uuid        NOT NULL REFERENCES patients (id),
    source_system text        NOT NULL,
    study_id      text        NOT NULL,
    modality      text        NOT NULL CHECK (modality IN ('CT', 'DX', 'MG')),
    body_site     text        NOT NULL DEFAULT '',
    performed_at  timestamptz NOT NULL,
    conclusion    text        NOT NULL,
    fingerprint   bytea       NOT NULL,
    status        text        NOT NULL CHECK (status IN ('received', 'assessed')),
    received_at   timestamptz NOT NULL,
    assessed_at   timestamptz,
    UNIQUE (source_system, study_id)
);

-- +goose Down
DROP TABLE intakes;
