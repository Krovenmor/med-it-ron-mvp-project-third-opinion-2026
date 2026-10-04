package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/infra/postgres/queries"
)

type Recommendations struct {
	executor
	q queries.Recommendations
}

func NewRecommendations(pool *pgxpool.Pool, q queries.Queries) *Recommendations {
	return &Recommendations{executor: executor{pool}, q: q.Recommendations}
}

func (r *Recommendations) InsertMany(ctx context.Context, recs []domain.Recommendation) error {
	batch := &pgx.Batch{}
	for _, rec := range recs {
		batch.Queue(r.q.Insert, insertRecommendationArgs(rec))
	}
	if err := r.db(ctx).SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("insert recommendations: %w", err)
	}
	return nil
}

func (r *Recommendations) ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.Recommendation, error) {
	rows, err := r.db(ctx).Query(ctx, r.q.ListByCase, pgx.StrictNamedArgs{"case_id": caseID})
	if err != nil {
		return nil, fmt.Errorf("list recommendations: %w", err)
	}
	found, err := pgx.CollectRows(rows, pgx.RowToStructByName[recommendationRow])
	if err != nil {
		return nil, fmt.Errorf("scan recommendations: %w", err)
	}
	recs := make([]domain.Recommendation, 0, len(found))
	for _, row := range found {
		recs = append(recs, row.toDomain())
	}
	return recs, nil
}

func (r *Recommendations) ActiveServices(ctx context.Context, patientID, excludeCaseID uuid.UUID) ([]domain.Service, error) {
	rows, err := r.db(ctx).Query(ctx, r.q.ListActiveServices,
		pgx.StrictNamedArgs{"patient_id": patientID, "exclude_case_id": excludeCaseID},
	)
	if err != nil {
		return nil, fmt.Errorf("list active services: %w", err)
	}
	found, err := pgx.CollectRows(rows, pgx.RowToStructByName[serviceRow])
	if err != nil {
		return nil, fmt.Errorf("scan active services: %w", err)
	}
	services := make([]domain.Service, 0, len(found))
	for _, row := range found {
		services = append(services, row.toDomain())
	}
	return services, nil
}
