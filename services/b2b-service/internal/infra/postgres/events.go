package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/queries"
)

type Events struct {
	executor
	q queries.CaseEvents
}

func NewEvents(pool *pgxpool.Pool, q queries.Queries) *Events {
	return &Events{executor: executor{pool}, q: q.CaseEvents}
}

func (r *Events) Append(ctx context.Context, e domain.CaseEvent) error {
	if _, err := r.db(ctx).Exec(ctx, r.q.Insert, insertCaseEventArgs(e)); err != nil {
		return fmt.Errorf("append case event: %w", err)
	}
	return nil
}
