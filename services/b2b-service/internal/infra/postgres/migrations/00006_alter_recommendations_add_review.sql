-- +goose Up
ALTER TABLE recommendations
    ADD COLUMN mark           text CHECK (mark IN ('critical', 'minor', 'rejected')),
    ADD COLUMN reject_reason  text CHECK (reject_reason IN ('contraindicated', 'other')),
    ADD COLUMN reject_comment text,
    ADD COLUMN marked_by      text,
    ADD COLUMN marked_at      timestamptz,
    ADD CONSTRAINT recommendations_reject_reason_only_when_rejected CHECK (reject_reason IS NULL OR mark = 'rejected'),
    ADD CONSTRAINT recommendations_rejected_requires_reason CHECK (mark <> 'rejected' OR reject_reason IS NOT NULL);

-- +goose Down
ALTER TABLE recommendations
    DROP CONSTRAINT recommendations_rejected_requires_reason,
    DROP CONSTRAINT recommendations_reject_reason_only_when_rejected,
    DROP COLUMN marked_at,
    DROP COLUMN marked_by,
    DROP COLUMN reject_comment,
    DROP COLUMN reject_reason,
    DROP COLUMN mark;
