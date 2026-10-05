package jobs

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
)

type jobRow struct {
	ID       int64           `db:"id"`
	Kind     string          `db:"kind"`
	Key      string          `db:"key"`
	Payload  json.RawMessage `db:"payload"`
	RunAt    time.Time       `db:"run_at"`
	Attempts int             `db:"attempts"`
}

func (r jobRow) toJob() Job {
	return Job{ID: r.ID, Kind: r.Kind, Key: r.Key, Payload: r.Payload, RunAt: r.RunAt, Attempts: r.Attempts}
}

func enqueueArgs(job Job, now time.Time) pgx.StrictNamedArgs {
	payload := job.Payload
	if payload == nil {
		payload = json.RawMessage("{}")
	}
	return pgx.StrictNamedArgs{
		"kind":       job.Kind,
		"key":        job.Key,
		"payload":    string(payload),
		"run_at":     job.RunAt,
		"created_at": now,
	}
}

func claimArgs(now time.Time, lease time.Duration) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{"now": now, "lease_seconds": lease.Seconds()}
}

func settleArgs(job Job) pgx.StrictNamedArgs {
	return pgx.StrictNamedArgs{"id": job.ID, "attempts": job.Attempts}
}

func rescheduleArgs(job Job, runAt time.Time, cause string) pgx.StrictNamedArgs {
	args := settleArgs(job)
	args["run_at"] = runAt
	args["last_error"] = cause
	return args
}

func failArgs(job Job, cause string) pgx.StrictNamedArgs {
	args := settleArgs(job)
	args["last_error"] = cause
	return args
}
