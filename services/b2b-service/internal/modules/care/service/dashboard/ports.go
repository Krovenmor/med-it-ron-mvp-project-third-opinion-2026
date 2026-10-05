package dashboard

import (
	"context"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
)

type Clock interface {
	Now() time.Time
}

type Facts interface {
	CaseFacts(ctx context.Context, from, to time.Time) ([]domain.CaseFact, error)
	TaskFacts(ctx context.Context, from, to time.Time) ([]domain.TaskFact, error)
	DeclineReasons(ctx context.Context, from, to time.Time) ([]domain.DeclineCount, error)
}
