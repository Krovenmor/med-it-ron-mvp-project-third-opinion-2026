package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/queries"
)

type Cases struct {
	executor
	q queries.Cases
}

func NewCases(pool *pgxpool.Pool, q queries.Queries) *Cases {
	return &Cases{executor: executor{pool}, q: q.Cases}
}

func (r *Cases) CreateIfAbsent(ctx context.Context, c domain.Case) (uuid.UUID, bool, error) {
	var id uuid.UUID
	err := r.db(ctx).QueryRow(ctx, r.q.InsertIfAbsent, insertCaseArgs(c)).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, false, nil
	}
	if err != nil {
		return uuid.Nil, false, fmt.Errorf("insert case: %w", err)
	}
	return id, true, nil
}

func (r *Cases) Get(ctx context.Context, id uuid.UUID) (domain.Case, error) {
	return r.getOne(ctx, r.q.Get, pgx.StrictNamedArgs{"id": id})
}

func (r *Cases) GetForUpdate(ctx context.Context, id uuid.UUID) (domain.Case, error) {
	return r.getOne(ctx, r.q.GetForUpdate, pgx.StrictNamedArgs{"id": id})
}

func (r *Cases) GetByStudy(ctx context.Context, sourceSystem, studyID string) (domain.Case, error) {
	return r.getOne(ctx, r.q.GetByStudy, pgx.StrictNamedArgs{"source_system": sourceSystem, "study_id": studyID})
}

func (r *Cases) Update(ctx context.Context, c domain.Case) error {
	tag, err := r.db(ctx).Exec(ctx, r.q.Update, updateCaseArgs(c))
	if err != nil {
		return fmt.Errorf("update case: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("case %s: %w", c.ID, domain.ErrNotFound)
	}
	return nil
}

func (r *Cases) ReviewQueue(ctx context.Context) ([]domain.ReviewQueueItem, error) {
	rows, err := r.db(ctx).Query(ctx, r.q.ListReviewQueue)
	if err != nil {
		return nil, fmt.Errorf("list review queue: %w", err)
	}
	found, err := pgx.CollectRows(rows, pgx.RowToStructByName[reviewQueueRow])
	if err != nil {
		return nil, fmt.Errorf("scan review queue: %w", err)
	}
	items := make([]domain.ReviewQueueItem, 0, len(found))
	for _, row := range found {
		items = append(items, row.toDomain())
	}
	return items, nil
}

func (r *Cases) ListByPatient(ctx context.Context, patientID uuid.UUID) ([]domain.Case, error) {
	rows, err := r.db(ctx).Query(ctx, r.q.ListByPatient, pgx.StrictNamedArgs{"patient_id": patientID})
	if err != nil {
		return nil, fmt.Errorf("list patient cases: %w", err)
	}
	found, err := pgx.CollectRows(rows, pgx.RowToStructByName[caseRow])
	if err != nil {
		return nil, fmt.Errorf("scan patient cases: %w", err)
	}
	cases := make([]domain.Case, 0, len(found))
	for _, row := range found {
		cases = append(cases, row.toDomain())
	}
	return cases, nil
}

func (r *Cases) getOne(ctx context.Context, q string, args pgx.StrictNamedArgs) (domain.Case, error) {
	rows, err := r.db(ctx).Query(ctx, q, args)
	if err != nil {
		return domain.Case{}, fmt.Errorf("get case: %w", err)
	}
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[caseRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Case{}, fmt.Errorf("case: %w", domain.ErrNotFound)
	}
	if err != nil {
		return domain.Case{}, fmt.Errorf("scan case: %w", err)
	}
	return row.toDomain(), nil
}
