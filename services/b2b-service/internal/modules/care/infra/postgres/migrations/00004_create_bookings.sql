-- +goose Up
CREATE TABLE bookings (
    id                uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id           uuid        NOT NULL REFERENCES routes (case_id),
    recommendation_id uuid        NOT NULL REFERENCES plan_items (id),
    appointment_id    text        NOT NULL UNIQUE,
    scheduled_at      timestamptz NOT NULL,
    channel           text        NOT NULL CHECK (channel IN ('self', 'operator')),
    created_at        timestamptz NOT NULL
);

CREATE INDEX bookings_case_id_idx ON bookings (case_id);

-- +goose Down
DROP TABLE bookings;
