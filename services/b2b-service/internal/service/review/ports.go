package review

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Clock interface {
	Now() time.Time
}

type Cases interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Case, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (domain.Case, error)
	Update(ctx context.Context, c domain.Case) error
	ReviewQueue(ctx context.Context) ([]domain.ReviewQueueItem, error)
}

type Recommendations interface {
	Get(ctx context.Context, caseID, id uuid.UUID) (domain.Recommendation, error)
	ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.Recommendation, error)
	UpdateReview(ctx context.Context, rec domain.Recommendation) error
	Append(ctx context.Context, rec domain.Recommendation) (domain.Recommendation, error)
}

type Events interface {
	Append(ctx context.Context, e domain.CaseEvent) error
}
