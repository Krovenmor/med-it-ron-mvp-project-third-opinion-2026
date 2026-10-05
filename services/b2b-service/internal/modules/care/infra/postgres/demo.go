package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Demo struct {
	postgres.Executor
	q queries.Queries
}

func NewDemo(db *postgres.Database, q queries.Queries) *Demo {
	return &Demo{Executor: db.Executor(), q: q}
}

func (r *Demo) Reset(ctx context.Context) error {
	if _, err := r.DB(ctx).Exec(ctx, r.q.Demo.Reset); err != nil {
		return fmt.Errorf("reset care data: %w", err)
	}
	return nil
}

func (r *Demo) Import(ctx context.Context, cases []history.Case) error {
	batch := &pgx.Batch{}
	for _, c := range cases {
		batch.Queue(r.q.Patients.Upsert, upsertPatientArgs(historyPatient(c), c.ReceivedAt))
		batch.Queue(r.q.Routes.InsertIfAbsent, routeArgs(historyRoute(c)))
		for _, item := range historyPlanItems(c) {
			batch.Queue(r.q.PlanItems.Insert, planItemArgs(item))
		}
		for _, b := range c.Bookings {
			batch.Queue(r.q.History.InsertBooking, historyBookingArgs(historyBooking(c, b)))
		}
		for _, t := range c.Tasks {
			batch.Queue(r.q.History.InsertTask, historyTaskArgs(historyTask(c, t)))
			for _, a := range t.Attempts {
				batch.Queue(r.q.CallAttempts.Insert, insertCallAttemptArgs(historyAttempt(t, a)))
			}
		}
		for _, e := range historyEvents(c) {
			batch.Queue(r.q.CaseEvents.Insert, insertCaseEventArgs(e))
		}
	}
	if err := r.DB(ctx).SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("import care history: %w", err)
	}
	return nil
}
