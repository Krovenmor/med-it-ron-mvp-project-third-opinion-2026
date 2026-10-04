package plan

import (
	"context"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Patients interface {
	GetByExternalID(ctx context.Context, sourceSystem, externalID string) (domain.Patient, error)
}

type Cases interface {
	ListByPatient(ctx context.Context, patientID uuid.UUID) ([]domain.Case, error)
}

type Recommendations interface {
	ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.Recommendation, error)
}
