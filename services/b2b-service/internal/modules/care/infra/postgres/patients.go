package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Patients struct {
	postgres.Executor
	q queries.Patients
}

func NewPatients(db *postgres.Database, q queries.Queries) *Patients {
	return &Patients{Executor: db.Executor(), q: q.Patients}
}

func (r *Patients) Upsert(ctx context.Context, p domain.Patient, now time.Time) error {
	if _, err := r.DB(ctx).Exec(ctx, r.q.Upsert, upsertPatientArgs(p, now)); err != nil {
		return fmt.Errorf("upsert patient: %w", err)
	}
	return nil
}

func (r *Patients) Get(ctx context.Context, id uuid.UUID) (domain.Patient, error) {
	return r.getOne(ctx, r.q.Get, pgx.StrictNamedArgs{"id": id})
}

func (r *Patients) GetByExternalID(ctx context.Context, sourceSystem, externalID string) (domain.Patient, error) {
	return r.getOne(ctx, r.q.GetByExternalID, pgx.StrictNamedArgs{"source_system": sourceSystem, "external_id": externalID})
}

func (r *Patients) getOne(ctx context.Context, q string, args pgx.StrictNamedArgs) (domain.Patient, error) {
	rows, err := r.DB(ctx).Query(ctx, q, args)
	if err != nil {
		return domain.Patient{}, fmt.Errorf("get patient: %w", err)
	}
	p, found, err := postgres.CollectOne(rows, patientRow.toDomain)
	switch {
	case err != nil:
		return domain.Patient{}, fmt.Errorf("scan patient: %w", err)
	case !found:
		return domain.Patient{}, fmt.Errorf("patient: %w", apperr.ErrNotFound)
	}
	return p, nil
}
