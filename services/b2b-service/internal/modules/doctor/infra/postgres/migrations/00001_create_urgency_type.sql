-- +goose Up
CREATE TYPE urgency AS ENUM ('normal', 'planned', 'priority', 'emergency');

-- +goose Down
DROP TYPE urgency;
