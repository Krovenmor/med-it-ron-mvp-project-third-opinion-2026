package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/cases"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/review"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
)

type Cases interface {
	Get(ctx context.Context, id uuid.UUID) (cases.Details, error)
	History(ctx context.Context, id uuid.UUID) (gatewayapi.History, error)
}

type Review interface {
	Queue(ctx context.Context) ([]domain.ReviewQueueItem, error)
	Open(ctx context.Context, caseID uuid.UUID, actor string) error
	ReviewRecommendation(ctx context.Context, cmd review.ReviewRecommendation) (domain.Recommendation, error)
	AddRecommendation(ctx context.Context, cmd review.AddRecommendation) (domain.Recommendation, error)
	ChangeUrgency(ctx context.Context, cmd review.ChangeUrgency) error
	Confirm(ctx context.Context, caseID uuid.UUID, actor string) (domain.Case, error)
}
