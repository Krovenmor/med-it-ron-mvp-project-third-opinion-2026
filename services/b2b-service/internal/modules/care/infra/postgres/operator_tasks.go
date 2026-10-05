package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type OperatorTasks struct {
	postgres.Executor
	q queries.OperatorTasks
}

func NewOperatorTasks(db *postgres.Database, q queries.Queries) *OperatorTasks {
	return &OperatorTasks{Executor: db.Executor(), q: q.OperatorTasks}
}

func (r *OperatorTasks) CreateIfNoActive(ctx context.Context, t domain.OperatorTask) (domain.OperatorTask, bool, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.InsertIfNoActive, insertTaskArgs(t))
	if err != nil {
		return domain.OperatorTask{}, false, fmt.Errorf("insert operator task: %w", err)
	}
	task, created, err := postgres.CollectOne(rows, taskRow.toDomain)
	if err != nil {
		return domain.OperatorTask{}, false, fmt.Errorf("scan operator task: %w", err)
	}
	return task, created, nil
}

func (r *OperatorTasks) Get(ctx context.Context, id uuid.UUID) (domain.OperatorTask, error) {
	return r.mustGet(ctx, r.q.Get, pgx.StrictNamedArgs{"id": id})
}

func (r *OperatorTasks) GetForUpdate(ctx context.Context, id uuid.UUID) (domain.OperatorTask, error) {
	return r.mustGet(ctx, r.q.GetForUpdate, pgx.StrictNamedArgs{"id": id})
}

func (r *OperatorTasks) ActiveByCaseForUpdate(ctx context.Context, caseID uuid.UUID) (domain.OperatorTask, bool, error) {
	return r.get(ctx, r.q.GetActiveByCaseForUpdate, pgx.StrictNamedArgs{"case_id": caseID})
}

func (r *OperatorTasks) Update(ctx context.Context, t domain.OperatorTask) error {
	tag, err := r.DB(ctx).Exec(ctx, r.q.Update, updateTaskArgs(t))
	if err != nil {
		return fmt.Errorf("update operator task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("operator task %s: %w", t.ID, apperr.ErrNotFound)
	}
	return nil
}

func (r *OperatorTasks) ListActive(ctx context.Context, now time.Time, assignee string) ([]domain.TaskListItem, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.ListActive, pgx.StrictNamedArgs{"now": now, "assignee": assignee})
	if err != nil {
		return nil, fmt.Errorf("list operator tasks: %w", err)
	}
	items, err := postgres.CollectAll(rows, taskListRow.toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan operator tasks: %w", err)
	}
	return items, nil
}

func (r *OperatorTasks) mustGet(ctx context.Context, q string, args pgx.StrictNamedArgs) (domain.OperatorTask, error) {
	task, found, err := r.get(ctx, q, args)
	switch {
	case err != nil:
		return domain.OperatorTask{}, err
	case !found:
		return domain.OperatorTask{}, fmt.Errorf("operator task: %w", apperr.ErrNotFound)
	}
	return task, nil
}

func (r *OperatorTasks) get(ctx context.Context, q string, args pgx.StrictNamedArgs) (domain.OperatorTask, bool, error) {
	rows, err := r.DB(ctx).Query(ctx, q, args)
	if err != nil {
		return domain.OperatorTask{}, false, fmt.Errorf("get operator task: %w", err)
	}
	task, found, err := postgres.CollectOne(rows, taskRow.toDomain)
	if err != nil {
		return domain.OperatorTask{}, false, fmt.Errorf("scan operator task: %w", err)
	}
	return task, found, nil
}
