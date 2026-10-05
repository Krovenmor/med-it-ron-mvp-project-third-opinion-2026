package events

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
)

type outgoingEnvelope struct {
	Topic      string    `json:"topic"`
	Data       any       `json:"data"`
	OccurredAt time.Time `json:"occurred_at"`
}

type incomingEnvelope struct {
	Topic      string          `json:"topic"`
	Data       json.RawMessage `json:"data"`
	OccurredAt time.Time       `json:"occurred_at"`
}

func rememberArgs(m Message, now time.Time) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{"message_id": m.ID, "topic": m.Topic, "processed_at": now}
}
