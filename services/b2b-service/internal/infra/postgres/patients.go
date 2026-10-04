package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/queries"
)

type Patients struct {
	executor
	q queries.Patients
}

func NewPatients(pool *pgxpool.Pool, q queries.Queries) *Patients {
	return &Patients{executor: executor{pool}, q: q.Patients}
}

func (r *Patients) Upsert(ctx context.Context, p domain.Patient, now time.Time) (uuid.UUID, error) {
	var id uuid.UUID
	if err := r.db(ctx).QueryRow(ctx, r.q.Upsert, upsertPatientArgs(p, now)).Scan(&id); err != nil {
		return uuid.Nil, fmt.Errorf("upsert patient: %w", err)
	}
	return id, nil
}

func (r *Patients) Get(ctx context.Context, id uuid.UUID) (domain.Patient, error) {
	return r.getOne(ctx, r.q.Get, pgx.StrictNamedArgs{"id": id})
}

func (r *Patients) GetByExternalID(ctx context.Context, sourceSystem, externalID string) (domain.Patient, error) {
	return r.getOne(ctx, r.q.GetByExternalID, pgx.StrictNamedArgs{"source_system": sourceSystem, "external_id": externalID})
}

func (r *Patients) getOne(ctx context.Context, q string, args pgx.StrictNamedArgs) (domain.Patient, error) {
	rows, err := r.db(ctx).Query(ctx, q, args)
	if err != nil {
		return domain.Patient{}, fmt.Errorf("get patient: %w", err)
	}
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[patientRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Patient{}, fmt.Errorf("patient: %w", domain.ErrNotFound)
	}
	if err != nil {
		return domain.Patient{}, fmt.Errorf("scan patient: %w", err)
	}
	return row.toDomain(), nil
}
