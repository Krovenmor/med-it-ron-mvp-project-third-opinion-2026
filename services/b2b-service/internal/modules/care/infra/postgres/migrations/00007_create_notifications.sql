-- +goose Up
CREATE TABLE notifications (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id    uuid        NOT NULL REFERENCES routes (case_id),
    recipient  text        NOT NULL CHECK (recipient IN ('patient', 'duty_doctor', 'admin_on_duty')),
    channel    text        NOT NULL CHECK (channel IN ('push', 'sms', 'staff')),
    kind       text        NOT NULL,
    text       text        NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE INDEX notifications_created_at_idx ON notifications (created_at DESC);

-- +goose Down
DROP TABLE notifications;
