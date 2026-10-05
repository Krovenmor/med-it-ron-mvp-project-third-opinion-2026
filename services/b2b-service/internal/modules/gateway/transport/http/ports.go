package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/service/intake"
)

type Intake interface {
	Ingest(ctx context.Context, report domain.Report) (intake.IngestResult, error)
	Status(ctx context.Context, caseID uuid.UUID) (domain.Intake, error)
}
