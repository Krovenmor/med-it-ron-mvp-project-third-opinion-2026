-- +goose Up
ALTER TABLE case_events
    ADD COLUMN payload jsonb NOT NULL DEFAULT '{}';

-- +goose Down
ALTER TABLE case_events
    DROP COLUMN payload;
