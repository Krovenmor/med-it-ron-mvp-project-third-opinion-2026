package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Dashboard struct {
	postgres.Executor
	q queries.Dashboard
}

func NewDashboard(db *postgres.Database, q queries.Queries) *Dashboard {
	return &Dashboard{Executor: db.Executor(), q: q.Dashboard}
}

func (r *Dashboard) CaseFacts(ctx context.Context, from, to time.Time) ([]domain.CaseFact, error) {
	return collectFacts(ctx, r, r.q.CaseFacts, from, to, caseFactRow.toDomain)
}

func (r *Dashboard) TaskFacts(ctx context.Context, from, to time.Time) ([]domain.TaskFact, error) {
	return collectFacts(ctx, r, r.q.TaskFacts, from, to, taskFactRow.toDomain)
}

func (r *Dashboard) DeclineReasons(ctx context.Context, from, to time.Time) ([]domain.DeclineCount, error) {
	return collectFacts(ctx, r, r.q.DeclineReasons, from, to, declineCountRow.toDomain)
}

func collectFacts[Row, T any](ctx context.Context, r *Dashboard, q string, from, to time.Time, toDomain func(Row) T) ([]T, error) {
	rows, err := r.DB(ctx).Query(ctx, q, pgx.StrictNamedArgs{"from": from, "to": to})
	if err != nil {
		return nil, fmt.Errorf("query dashboard: %w", err)
	}
	items, err := postgres.CollectAll(rows, toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan dashboard: %w", err)
	}
	return items, nil
}
