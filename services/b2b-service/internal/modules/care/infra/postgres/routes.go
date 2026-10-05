package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Routes struct {
	postgres.Executor
	q queries.Routes
}

func NewRoutes(db *postgres.Database, q queries.Queries) *Routes {
	return &Routes{Executor: db.Executor(), q: q.Routes}
}

func (r *Routes) CreateIfAbsent(ctx context.Context, route domain.Route) (bool, error) {
	tag, err := r.DB(ctx).Exec(ctx, r.q.InsertIfAbsent, routeArgs(route))
	if err != nil {
		return false, fmt.Errorf("insert route: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

func (r *Routes) Get(ctx context.Context, caseID uuid.UUID) (domain.Route, error) {
	return r.getOne(ctx, r.q.Get, caseID)
}

func (r *Routes) GetForUpdate(ctx context.Context, caseID uuid.UUID) (domain.Route, error) {
	return r.getOne(ctx, r.q.GetForUpdate, caseID)
}

func (r *Routes) Update(ctx context.Context, route domain.Route) error {
	tag, err := r.DB(ctx).Exec(ctx, r.q.Update, routeArgs(route))
	if err != nil {
		return fmt.Errorf("update route: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("route %s: %w", route.CaseID, apperr.ErrNotFound)
	}
	return nil
}

func (r *Routes) ListByPatient(ctx context.Context, patientID uuid.UUID) ([]domain.Route, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.ListByPatient, pgx.StrictNamedArgs{"patient_id": patientID})
	if err != nil {
		return nil, fmt.Errorf("list patient routes: %w", err)
	}
	routes, err := postgres.CollectAll(rows, routeRow.toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan patient routes: %w", err)
	}
	return routes, nil
}

func (r *Routes) getOne(ctx context.Context, q string, caseID uuid.UUID) (domain.Route, error) {
	rows, err := r.DB(ctx).Query(ctx, q, pgx.StrictNamedArgs{"case_id": caseID})
	if err != nil {
		return domain.Route{}, fmt.Errorf("get route: %w", err)
	}
	route, found, err := postgres.CollectOne(rows, routeRow.toDomain)
	switch {
	case err != nil:
		return domain.Route{}, fmt.Errorf("scan route: %w", err)
	case !found:
		return domain.Route{}, fmt.Errorf("case %s: %w", caseID, apperr.ErrNotFound)
	}
	return route, nil
}
