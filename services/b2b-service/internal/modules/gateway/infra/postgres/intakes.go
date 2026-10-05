package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Intakes struct {
	postgres.Executor
	q queries.Intakes
}

func NewIntakes(db *postgres.Database, q queries.Queries) *Intakes {
	return &Intakes{Executor: db.Executor(), q: q.Intakes}
}

func (r *Intakes) CreateIfAbsent(ctx context.Context, i domain.Intake) (uuid.UUID, bool, error) {
	var id uuid.UUID
	err := r.DB(ctx).QueryRow(ctx, r.q.InsertIfAbsent, insertIntakeArgs(i)).Scan(&id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return uuid.Nil, false, nil
	case err != nil:
		return uuid.Nil, false, fmt.Errorf("insert intake: %w", err)
	}
	return id, true, nil
}

func (r *Intakes) Get(ctx context.Context, caseID uuid.UUID) (domain.Intake, error) {
	return r.getOne(ctx, r.q.Get, pgx.StrictNamedArgs{"case_id": caseID})
}

func (r *Intakes) GetForUpdate(ctx context.Context, caseID uuid.UUID) (domain.Intake, error) {
	return r.getOne(ctx, r.q.GetForUpdate, pgx.StrictNamedArgs{"case_id": caseID})
}

func (r *Intakes) GetByStudy(ctx context.Context, sourceSystem, studyID string) (domain.Intake, error) {
	return r.getOne(ctx, r.q.GetByStudy, pgx.StrictNamedArgs{"source_system": sourceSystem, "study_id": studyID})
}

func (r *Intakes) Update(ctx context.Context, i domain.Intake) error {
	tag, err := r.DB(ctx).Exec(ctx, r.q.Update, updateIntakeArgs(i))
	if err != nil {
		return fmt.Errorf("update intake: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("intake %s: %w", i.CaseID, apperr.ErrNotFound)
	}
	return nil
}

func (r *Intakes) getOne(ctx context.Context, q string, args pgx.StrictNamedArgs) (domain.Intake, error) {
	rows, err := r.DB(ctx).Query(ctx, q, args)
	if err != nil {
		return domain.Intake{}, fmt.Errorf("get intake: %w", err)
	}
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[intakeRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Intake{}, fmt.Errorf("intake: %w", apperr.ErrNotFound)
	}
	if err != nil {
		return domain.Intake{}, fmt.Errorf("scan intake: %w", err)
	}
	return row.toDomain(), nil
}
