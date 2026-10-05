package plan

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
)

type Clock interface {
	Now() time.Time
}

type Patients interface {
	GetByExternalID(ctx context.Context, sourceSystem, externalID string) (domain.Patient, error)
}

type Routes interface {
	ListByPatient(ctx context.Context, patientID uuid.UUID) ([]domain.Route, error)
}

type PlanItems interface {
	ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.PlanItem, error)
}
