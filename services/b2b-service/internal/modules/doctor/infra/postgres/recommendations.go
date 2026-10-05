package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Recommendations struct {
	postgres.Executor
	q queries.Recommendations
}

func NewRecommendations(db *postgres.Database, q queries.Queries) *Recommendations {
	return &Recommendations{Executor: db.Executor(), q: q.Recommendations}
}

func (r *Recommendations) InsertMany(ctx context.Context, recs []domain.Recommendation) error {
	batch := &pgx.Batch{}
	for _, rec := range recs {
		batch.Queue(r.q.Insert, insertRecommendationArgs(rec))
	}
	if err := r.DB(ctx).SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("insert recommendations: %w", err)
	}
	return nil
}

func (r *Recommendations) Append(ctx context.Context, rec domain.Recommendation) (domain.Recommendation, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.Append, appendRecommendationArgs(rec))
	if err != nil {
		return domain.Recommendation{}, fmt.Errorf("append recommendation: %w", err)
	}
	saved, _, err := postgres.CollectOne(rows, recommendationRow.toDomain)
	if err != nil {
		return domain.Recommendation{}, fmt.Errorf("scan appended recommendation: %w", err)
	}
	return saved, nil
}

func (r *Recommendations) Get(ctx context.Context, caseID, id uuid.UUID) (domain.Recommendation, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.Get, pgx.StrictNamedArgs{"id": id, "case_id": caseID})
	if err != nil {
		return domain.Recommendation{}, fmt.Errorf("get recommendation: %w", err)
	}
	rec, found, err := postgres.CollectOne(rows, recommendationRow.toDomain)
	switch {
	case err != nil:
		return domain.Recommendation{}, fmt.Errorf("scan recommendation: %w", err)
	case !found:
		return domain.Recommendation{}, fmt.Errorf("recommendation %s: %w", id, apperr.ErrNotFound)
	}
	return rec, nil
}

func (r *Recommendations) UpdateReview(ctx context.Context, rec domain.Recommendation) error {
	return r.exec(ctx, r.q.UpdateReview, updateReviewArgs(rec), rec.ID)
}

func (r *Recommendations) UpdatePatientText(ctx context.Context, rec domain.Recommendation) error {
	return r.exec(ctx, r.q.UpdatePatientText, pgx.StrictNamedArgs{"id": rec.ID, "patient_text": rec.PatientText}, rec.ID)
}

func (r *Recommendations) ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.Recommendation, error) {
	rows, err := r.DB(ctx).Query(ctx, r.q.ListByCase, pgx.StrictNamedArgs{"case_id": caseID})
	if err != nil {
		return nil, fmt.Errorf("list recommendations: %w", err)
	}
	recs, err := postgres.CollectAll(rows, recommendationRow.toDomain)
	if err != nil {
		return nil, fmt.Errorf("scan recommendations: %w", err)
	}
	return recs, nil
}

func (r *Recommendations) exec(ctx context.Context, q string, args pgx.StrictNamedArgs, id uuid.UUID) error {
	tag, err := r.DB(ctx).Exec(ctx, q, args)
	if err != nil {
		return fmt.Errorf("update recommendation: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("recommendation %s: %w", id, apperr.ErrNotFound)
	}
	return nil
}
