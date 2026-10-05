package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Cases struct {
	postgres.Executor
	q queries.Cases
}

func NewCases(db *postgres.Database, q queries.Queries) *Cases {
	return &Cases{Executor: db.Executor(), q: q.Cases}
}

func (r *Cases) CreateIfAbsent(ctx context.Context, c domain.Case) (bool, error) {
	tag, err := r.DB(ctx).Exec(ctx, r.q.InsertIfAbsent, insertCaseArgs(c))
	if err != nil {
		return false, fmt.Errorf("insert case: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

func (r *Cases) Get(ctx context.Context, id uuid.UUID) (domain.Case, error) {
	return r.getOne(ctx, r.q.Get, id)
}

func (r *Cases) GetForUpdate(ctx context.Context, id uuid.UUID) (domain.Case, error) {
	return r.getOne(ctx, r.q.GetForUpdate, id)
}

func (r *Cases) Update(ctx context.Context, c domain.Case) error {
	tag, err := r.DB(ctx).Exec(ctx, r.q.Update, updateCaseArgs(c))
	if err != nil {
		return fmt.Errorf("update case: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("case %s: %w", c.ID, apperr.ErrNotFound)
	}
	return nil
}

func (r *Cases) ReviewQueue(ctx context.Context) ([]domain.ReviewQueueItem, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.ListReviewQueue)
	if err != nil {
		return nil, fmt.Errorf("list review queue: %w", err)
	}
	items, err := postgres.CollectAll(rows, reviewQueueRow.toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan review queue: %w", err)
	}
	return items, nil
}

func (r *Cases) getOne(ctx context.Context, q string, id uuid.UUID) (domain.Case, error) {
	rows, err := r.DB(ctx).Query(ctx, q, pgx.StrictNamedArgs{"id": id})
	if err != nil {
		return domain.Case{}, fmt.Errorf("get case: %w", err)
	}
	c, found, err := postgres.CollectOne(rows, caseRow.toDomain)
	switch {
	case err != nil:
		return domain.Case{}, fmt.Errorf("scan case: %w", err)
	case !found:
		return domain.Case{}, fmt.Errorf("case %s: %w", id, apperr.ErrNotFound)
	}
	return c, nil
}
