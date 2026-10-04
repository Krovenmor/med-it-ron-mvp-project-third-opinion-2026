package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/queries"
)

type Jobs struct {
	executor
	q queries.Jobs
}

func NewJobs(pool *pgxpool.Pool, q queries.Queries) *Jobs {
	return &Jobs{executor: executor{pool}, q: q.Jobs}
}

func (r *Jobs) Enqueue(ctx context.Context, job domain.Job, now time.Time) error {
	if _, err := r.db(ctx).Exec(ctx, r.q.Enqueue, enqueueJobArgs(job, now)); err != nil {
		return fmt.Errorf("enqueue job: %w", err)
	}
	return nil
}

func (r *Jobs) Claim(ctx context.Context, now time.Time, lease time.Duration) (domain.Job, bool, error) {
	rows, err := r.db(ctx).Query(ctx, r.q.Claim, claimJobArgs(now, lease))
	if err != nil {
		return domain.Job{}, false, fmt.Errorf("claim job: %w", err)
	}
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[jobRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Job{}, false, nil
	}
	if err != nil {
		return domain.Job{}, false, fmt.Errorf("scan claimed job: %w", err)
	}
	return row.toDomain(), true, nil
}

func (r *Jobs) Complete(ctx context.Context, job domain.Job) error {
	return r.settle(ctx, r.q.Complete, completeJobArgs(job))
}

func (r *Jobs) Reschedule(ctx context.Context, job domain.Job, runAt time.Time, cause string) error {
	return r.settle(ctx, r.q.Reschedule, rescheduleJobArgs(job, runAt, cause))
}

func (r *Jobs) Fail(ctx context.Context, job domain.Job, cause string) error {
	return r.settle(ctx, r.q.Fail, failJobArgs(job, cause))
}

func (r *Jobs) settle(ctx context.Context, q string, args pgx.StrictNamedArgs) error {
	tag, err := r.db(ctx).Exec(ctx, q, args)
	if err != nil {
		return fmt.Errorf("settle job: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrLeaseLost
	}
	return nil
}
