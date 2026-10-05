package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type CallAttempts struct {
	postgres.Executor
	q queries.CallAttempts
}

func NewCallAttempts(db *postgres.Database, q queries.Queries) *CallAttempts {
	return &CallAttempts{Executor: db.Executor(), q: q.CallAttempts}
}

func (r *CallAttempts) Insert(ctx context.Context, a domain.CallAttempt) error {
	if _, err := r.DB(ctx).Exec(ctx, r.q.Insert, insertCallAttemptArgs(a)); err != nil {
		return fmt.Errorf("insert call attempt: %w", err)
	}
	return nil
}

func (r *CallAttempts) ListByTask(ctx context.Context, taskID uuid.UUID) ([]domain.CallAttempt, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.ListByTask, pgx.StrictNamedArgs{"task_id": taskID})
	if err != nil {
		return nil, fmt.Errorf("list call attempts: %w", err)
	}
	attempts, err := postgres.CollectAll(rows, callAttemptRow.toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan call attempts: %w", err)
	}
	return attempts, nil
}
