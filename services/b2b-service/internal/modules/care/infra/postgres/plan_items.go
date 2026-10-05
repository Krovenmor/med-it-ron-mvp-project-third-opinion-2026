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

type PlanItems struct {
	postgres.Executor
	q queries.PlanItems
}

func NewPlanItems(db *postgres.Database, q queries.Queries) *PlanItems {
	return &PlanItems{Executor: db.Executor(), q: q.PlanItems}
}

func (r *PlanItems) InsertMany(ctx context.Context, items []domain.PlanItem) error {
	batch := &pgx.Batch{}
	for _, item := range items {
		batch.Queue(r.q.Insert, planItemArgs(item))
	}
	if err := r.DB(ctx).SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("insert plan items: %w", err)
	}
	return nil
}

func (r *PlanItems) Get(ctx context.Context, caseID, id uuid.UUID) (domain.PlanItem, error) {
	return r.getOne(ctx, r.q.Get, caseID, id)
}

func (r *PlanItems) GetForUpdate(ctx context.Context, caseID, id uuid.UUID) (domain.PlanItem, error) {
	return r.getOne(ctx, r.q.GetForUpdate, caseID, id)
}

func (r *PlanItems) UpdateDecline(ctx context.Context, item domain.PlanItem) error {
	tag, err := r.DB(ctx).Exec(ctx, r.q.UpdateDecline, declineArgs(item))
	if err != nil {
		return fmt.Errorf("update plan item decline: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("recommendation %s: %w", item.ID, apperr.ErrNotFound)
	}
	return nil
}

func (r *PlanItems) getOne(ctx context.Context, q string, caseID, id uuid.UUID) (domain.PlanItem, error) {
	rows, err := r.DB(ctx).Query(ctx, q, pgx.StrictNamedArgs{"id": id, "case_id": caseID})
	if err != nil {
		return domain.PlanItem{}, fmt.Errorf("get plan item: %w", err)
	}
	item, found, err := postgres.CollectOne(rows, planItemRow.toDomain)
	switch {
	case err != nil:
		return domain.PlanItem{}, fmt.Errorf("scan plan item: %w", err)
	case !found:
		return domain.PlanItem{}, fmt.Errorf("recommendation %s: %w", id, apperr.ErrNotFound)
	}
	return item, nil
}

func (r *PlanItems) ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.PlanItem, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.ListByCase, pgx.StrictNamedArgs{"case_id": caseID})
	if err != nil {
		return nil, fmt.Errorf("list plan items: %w", err)
	}
	items, err := postgres.CollectAll(rows, planItemRow.toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan plan items: %w", err)
	}
	return items, nil
}

func (r *PlanItems) ActiveServices(ctx context.Context, patientID, excludeCaseID uuid.UUID) ([]domain.Service, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.ListActiveServices,
		pgx.StrictNamedArgs{"patient_id": patientID, "exclude_case_id": excludeCaseID},
	)
	if err != nil {
		return nil, fmt.Errorf("list active services: %w", err)
	}
	services, err := postgres.CollectAll(rows, serviceRow.toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan active services: %w", err)
	}
	return services, nil
}
