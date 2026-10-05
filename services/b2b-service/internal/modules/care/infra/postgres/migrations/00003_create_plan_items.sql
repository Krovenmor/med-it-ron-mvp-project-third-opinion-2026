-- +goose Up
CREATE TABLE plan_items (
    id           uuid PRIMARY KEY,
    case_id      uuid NOT NULL REFERENCES routes (case_id),
    position     int  NOT NULL,
    service_code text NOT NULL DEFAULT '',
    service_name text NOT NULL,
    patient_text text NOT NULL,
    mark         text NOT NULL CHECK (mark IN ('critical', 'minor'))
);

CREATE INDEX plan_items_case_id_idx ON plan_items (case_id, position);

-- +goose Down
DROP TABLE plan_items;
