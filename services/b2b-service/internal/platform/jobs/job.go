package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrLeaseLost = errors.New("job lease lost")

type Job struct {
	ID       int64
	Kind     string
	Key      string
	Payload  json.RawMessage
	RunAt    time.Time
	Attempts int
}

func New(kind, key string, runAt time.Time) Job {
	return Job{Kind: kind, Key: key, RunAt: runAt}
}

func NewWithPayload(kind, key string, runAt time.Time, payload any) (Job, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Job{}, fmt.Errorf("encode %s job payload: %w", kind, err)
	}
	job := New(kind, key, runAt)
	job.Payload = data
	return job, nil
}

func (j Job) Decode(dst any) error {
	if err := json.Unmarshal(j.Payload, dst); err != nil {
		return fmt.Errorf("decode %s job payload: %w", j.Kind, err)
	}
	return nil
}

type Config struct {
	Count        int           `env:"COUNT"`
	PollInterval time.Duration `env:"POLL_INTERVAL"`
	Lease        time.Duration `env:"LEASE"`
	JobTimeout   time.Duration `env:"JOB_TIMEOUT"`
	MaxAttempts  int           `env:"MAX_ATTEMPTS"`
	BackoffBase  time.Duration `env:"BACKOFF_BASE"`
	BackoffMax   time.Duration `env:"BACKOFF_MAX"`
}
