-- +goose Up
CREATE TABLE cases (
    id                    uuid        PRIMARY KEY,
    patient_id            uuid        NOT NULL,
    patient_source_system text        NOT NULL,
    patient_external_id   text        NOT NULL,
    patient_birth_date    date        NOT NULL,
    patient_sex           text        NOT NULL CHECK (patient_sex IN ('male', 'female')),
    study_id              text        NOT NULL,
    modality              text        NOT NULL CHECK (modality IN ('CT', 'DX', 'MG')),
    body_site             text        NOT NULL DEFAULT '',
    performed_at          timestamptz NOT NULL,
    conclusion            text        NOT NULL,
    status                text        NOT NULL CHECK (status IN ('in_review', 'confirmed')),
    urgency               urgency     NOT NULL,
    catalog_version       text        NOT NULL,
    guidelines_version    text        NOT NULL,
    received_at           timestamptz NOT NULL,
    assessed_at           timestamptz NOT NULL,
    confirmed_at          timestamptz,
    updated_at            timestamptz NOT NULL
);

CREATE INDEX cases_in_review_idx ON cases (urgency DESC, received_at) WHERE status = 'in_review';

-- +goose Down
DROP TABLE cases;
