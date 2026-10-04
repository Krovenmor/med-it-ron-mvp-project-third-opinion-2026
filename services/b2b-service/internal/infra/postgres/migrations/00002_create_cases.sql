-- +goose Up
CREATE TYPE urgency AS ENUM ('normal', 'planned', 'priority', 'emergency');

CREATE TABLE cases (
    id                 uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id         uuid        NOT NULL REFERENCES patients (id),
    source_system      text        NOT NULL,
    study_id           text        NOT NULL,
    modality           text        NOT NULL CHECK (modality IN ('CT', 'DX', 'MG')),
    body_site          text        NOT NULL DEFAULT '',
    performed_at       timestamptz NOT NULL,
    conclusion         text        NOT NULL,
    fingerprint        bytea       NOT NULL,
    status             text        NOT NULL CHECK (status IN ('draft', 'in_review', 'confirmed', 'notified', 'booked', 'completed', 'declined', 'unreachable')),
    urgency            urgency,
    catalog_version    text        NOT NULL DEFAULT '',
    guidelines_version text        NOT NULL DEFAULT '',
    received_at        timestamptz NOT NULL,
    updated_at         timestamptz NOT NULL,
    UNIQUE (source_system, study_id)
);

CREATE INDEX cases_patient_id_idx ON cases (patient_id);

-- +goose Down
DROP TABLE cases;
DROP TYPE urgency;
