package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Queue struct {
	postgres.Executor
	q queries.Queries
}

func NewQueue(db *postgres.Database) (*Queue, error) {
	q, err := queries.Load()
	if err != nil {
		return nil, err
	}
	return &Queue{Executor: db.Executor(), q: q}, nil
}

func (r *Queue) Enqueue(ctx context.Context, job Job, now time.Time) error {
	if _, err := r.DB(ctx).Exec(ctx, r.q.Enqueue, enqueueArgs(job, now)); err != nil {
		return fmt.Errorf("enqueue %s job: %w", job.Kind, err)
	}
	return nil
}

func (r *Queue) Claim(ctx context.Context, now time.Time, lease time.Duration) (Job, bool, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.Claim, claimArgs(now, lease))
	if err != nil {
		return Job{}, false, fmt.Errorf("claim job: %w", err)
	}
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[jobRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return Job{}, false, nil
	}
	if err != nil {
		return Job{}, false, fmt.Errorf("scan claimed job: %w", err)
	}
	return row.toJob(), true, nil
}

func (r *Queue) Complete(ctx context.Context, job Job) error {
	return r.settle(ctx, r.q.Complete, settleArgs(job))
}

func (r *Queue) Reschedule(ctx context.Context, job Job, runAt time.Time, cause string) error {
	return r.settle(ctx, r.q.Reschedule, rescheduleArgs(job, runAt, cause))
}

func (r *Queue) Fail(ctx context.Context, job Job, cause string) error {
	return r.settle(ctx, r.q.Fail, failArgs(job, cause))
}

func (r *Queue) CancelPending(ctx context.Context, key string, kinds []string) error {
	if _, err := r.DB(ctx).Exec(ctx, r.q.CancelPending, pgx.StrictNamedArgs{"key": key, "kinds": kinds}); err != nil {
		return fmt.Errorf("cancel pending jobs: %w", err)
	}
	return nil
}

func (r *Queue) settle(ctx context.Context, q string, args pgx.StrictNamedArgs) error {
	tag, err := r.DB(ctx).Exec(ctx, q, args)
	if err != nil {
		return fmt.Errorf("settle job: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrLeaseLost
	}
	return nil
}
