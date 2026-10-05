-- +goose Up
ALTER TABLE operator_tasks
    DROP CONSTRAINT operator_tasks_reason_check,
    ADD CONSTRAINT operator_tasks_reason_check CHECK (reason IN ('emergency', 'no_booking', 'help_request'));

-- +goose Down
ALTER TABLE operator_tasks
    DROP CONSTRAINT operator_tasks_reason_check,
    ADD CONSTRAINT operator_tasks_reason_check CHECK (reason IN ('emergency', 'no_booking'));
