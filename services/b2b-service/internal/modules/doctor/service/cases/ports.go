package cases

import (
	"context"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
)

type Cases interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Case, error)
}

type Recommendations interface {
	ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.Recommendation, error)
}

type MIS interface {
	PatientHistory(ctx context.Context, patientID uuid.UUID) (gatewayapi.History, error)
}
