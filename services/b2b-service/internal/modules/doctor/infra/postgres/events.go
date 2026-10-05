package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/infra/postgres/queries"
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

func (r *Events) Exists(ctx context.Context, caseID uuid.UUID, eventType domain.CaseEventType) (bool, error) {
	var exists bool
	err := r.DB(ctx).QueryRow(ctx, r.q.Exists, pgx.StrictNamedArgs{"case_id": caseID, "type": string(eventType)}).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check case event: %w", err)
	}
	return exists, nil
}
