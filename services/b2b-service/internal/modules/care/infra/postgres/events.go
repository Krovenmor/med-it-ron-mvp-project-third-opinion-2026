package postgres

import (
	"context"
	"fmt"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Events struct {
	postgres.Executor
	q queries.CaseEvents
}

func NewEvents(db *postgres.Database, q queries.Queries) *Events {
	return &Events{Executor: db.Executor(), q: q.CaseEvents}
}

func (r *Events) Append(ctx context.Context, e domain.CaseEvent) error {
	if _, err := r.DB(ctx).Exec(ctx, r.q.Insert, insertCaseEventArgs(e)); err != nil {
		return fmt.Errorf("append case event: %w", err)
	}
	return nil
}
