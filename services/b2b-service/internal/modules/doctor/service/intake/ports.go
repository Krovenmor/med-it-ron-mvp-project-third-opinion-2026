package intake

import (
	"context"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

type Cases interface {
	CreateIfAbsent(ctx context.Context, c domain.Case) (bool, error)
}

type Recommendations interface {
	InsertMany(ctx context.Context, recs []domain.Recommendation) error
}

type Events interface {
	Append(ctx context.Context, e domain.CaseEvent) error
}

type Jobs interface {
	Enqueue(ctx context.Context, job jobs.Job, now time.Time) error
}

type Publisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
}
