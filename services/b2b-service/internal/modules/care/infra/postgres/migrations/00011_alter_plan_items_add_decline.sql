-- +goose Up
ALTER TABLE plan_items
    ADD COLUMN declined_at     timestamptz,
    ADD COLUMN decline_reason  text CHECK (decline_reason IN ('expensive', 'far', 'other_clinic', 'not_needed', 'other')),
    ADD COLUMN decline_comment text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE plan_items
    DROP COLUMN decline_comment,
    DROP COLUMN decline_reason,
    DROP COLUMN declined_at;
