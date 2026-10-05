package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Demo struct {
	postgres.Executor
	q       queries.Demo
	history queries.History
	events  queries.CaseEvents
}

func NewDemo(db *postgres.Database, q queries.Queries) *Demo {
	return &Demo{Executor: db.Executor(), q: q.Demo, history: q.History, events: q.CaseEvents}
}

func (r *Demo) Reset(ctx context.Context) error {
	if _, err := r.DB(ctx).Exec(ctx, r.q.Reset); err != nil {
		return fmt.Errorf("reset doctor data: %w", err)
	}
	return nil
}

func (r *Demo) Import(ctx context.Context, cases []history.Case) error {
	batch := &pgx.Batch{}
	for _, c := range cases {
		batch.Queue(r.history.InsertCase, historyCaseArgs(c))
		for _, rec := range c.Recommendations {
			batch.Queue(r.history.InsertRecommendation, historyRecommendationArgs(c, rec))
		}
		for _, e := range historyEvents(c) {
			batch.Queue(r.events.Insert, insertCaseEventArgs(e))
		}
	}
	if err := r.DB(ctx).SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("import doctor history: %w", err)
	}
	return nil
}
