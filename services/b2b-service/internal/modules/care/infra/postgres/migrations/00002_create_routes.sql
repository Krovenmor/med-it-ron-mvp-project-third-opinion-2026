-- +goose Up
CREATE TABLE routes (
    case_id       uuid        PRIMARY KEY,
    patient_id    uuid        NOT NULL,
    status        text        NOT NULL CHECK (status IN ('received', 'in_review', 'confirmed', 'notified', 'booked', 'completed', 'declined', 'unreachable')),
    urgency       text        CHECK (urgency IN ('normal', 'planned', 'priority', 'emergency')),
    modality      text        NOT NULL CHECK (modality IN ('CT', 'DX', 'MG')),
    performed_at  timestamptz NOT NULL,
    received_at   timestamptz NOT NULL,
    assessed_at   timestamptz,
    review_due_at timestamptz,
    confirmed_at  timestamptz,
    updated_at    timestamptz NOT NULL
);

CREATE INDEX routes_patient_id_idx ON routes (patient_id);
CREATE INDEX routes_received_at_idx ON routes (received_at);

-- +goose Down
DROP TABLE routes;
