package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Demo struct {
	postgres.Executor
	q       queries.Demo
	history queries.History
}

func NewDemo(db *postgres.Database, q queries.Queries) *Demo {
	return &Demo{Executor: db.Executor(), q: q.Demo, history: q.History}
}

func (r *Demo) Reset(ctx context.Context) error {
	if _, err := r.DB(ctx).Exec(ctx, r.q.Reset); err != nil {
		return fmt.Errorf("reset gateway data: %w", err)
	}
	return nil
}

func (r *Demo) Import(ctx context.Context, cases []history.Case) error {
	batch := &pgx.Batch{}
	for _, c := range cases {
		batch.Queue(r.history.InsertPatient, historyPatientArgs(c))
		batch.Queue(r.history.InsertIntake, historyIntakeArgs(c))
	}
	if err := r.DB(ctx).SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("import gateway history: %w", err)
	}
	return nil
}
