-- +goose Up
CREATE TABLE operator_tasks (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id      uuid        NOT NULL REFERENCES routes (case_id),
    reason       text        NOT NULL CHECK (reason IN ('emergency', 'no_booking')),
    status       text        NOT NULL CHECK (status IN ('new', 'in_progress', 'no_answer', 'callback', 'contacted', 'booked', 'declined', 'handed_to_doctor', 'cancelled')),
    assignee     text,
    due_at       timestamptz NOT NULL,
    next_call_at timestamptz,
    attempts     int         NOT NULL DEFAULT 0,
    created_at   timestamptz NOT NULL,
    updated_at   timestamptz NOT NULL,
    closed_at    timestamptz
);

CREATE UNIQUE INDEX operator_tasks_one_active_per_case_idx ON operator_tasks (case_id)
    WHERE status IN ('new', 'in_progress', 'no_answer', 'callback');

-- +goose Down
DROP TABLE operator_tasks;
