package cases

import (
	"context"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Cases interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Case, error)
}

type Patients interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Patient, error)
}

type Recommendations interface {
	ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.Recommendation, error)
}
