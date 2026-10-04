package http

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/cases"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/intake"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/review"
)

type Intake interface {
	Ingest(ctx context.Context, report domain.Report) (intake.IngestResult, error)
}

type Cases interface {
	Get(ctx context.Context, id uuid.UUID) (cases.Details, error)
}

type Review interface {
	Queue(ctx context.Context) ([]domain.ReviewQueueItem, error)
	Open(ctx context.Context, caseID uuid.UUID, actor string) error
	ReviewRecommendation(ctx context.Context, cmd review.ReviewRecommendation) (domain.Recommendation, error)
	AddRecommendation(ctx context.Context, cmd review.AddRecommendation) (domain.Recommendation, error)
	ChangeUrgency(ctx context.Context, cmd review.ChangeUrgency) error
	Confirm(ctx context.Context, caseID uuid.UUID, actor string) (domain.Case, error)
}

type DemoClock interface {
	Advance(d time.Duration) time.Time
}
