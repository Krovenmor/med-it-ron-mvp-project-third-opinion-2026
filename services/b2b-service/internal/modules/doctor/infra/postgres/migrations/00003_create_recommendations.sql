-- +goose Up
CREATE TABLE recommendations (
    id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id        uuid        NOT NULL REFERENCES cases (id),
    position       int         NOT NULL,
    source         text        NOT NULL CHECK (source IN ('ai', 'doctor')),
    service_code   text        NOT NULL DEFAULT '',
    service_name   text        NOT NULL,
    importance     text        NOT NULL CHECK (importance IN ('high', 'low')),
    rationale      text        NOT NULL DEFAULT '',
    guideline_ref  text        NOT NULL DEFAULT '',
    patient_text   text        NOT NULL,
    already_booked boolean     NOT NULL DEFAULT false,
    mark           text        CHECK (mark IN ('critical', 'minor', 'rejected')),
    reject_reason  text        CHECK (reject_reason IN ('contraindicated', 'other')),
    reject_comment text,
    marked_by      text,
    marked_at      timestamptz,
    created_at     timestamptz NOT NULL,
    UNIQUE (case_id, position),
    CONSTRAINT recommendations_reject_reason_only_when_rejected CHECK (reject_reason IS NULL OR mark = 'rejected')
);

-- +goose Down
DROP TABLE recommendations;
