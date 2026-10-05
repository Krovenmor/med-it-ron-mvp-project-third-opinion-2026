-- +goose Up
CREATE TABLE operator_task_attempts (
    id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id        uuid        NOT NULL REFERENCES operator_tasks (id),
    outcome        text        NOT NULL CHECK (outcome IN ('no_answer', 'callback', 'contacted', 'declined', 'booked', 'handed_to_doctor')),
    decline_reason text        CHECK (decline_reason IN ('expensive', 'far', 'other_clinic', 'not_needed', 'other')),
    comment        text        NOT NULL DEFAULT '',
    callback_at    timestamptz,
    actor          text        NOT NULL,
    created_at     timestamptz NOT NULL
);

CREATE INDEX operator_task_attempts_task_id_idx ON operator_task_attempts (task_id, created_at);

-- +goose Down
DROP TABLE operator_task_attempts;
