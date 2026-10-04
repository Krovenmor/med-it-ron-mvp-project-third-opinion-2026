package worker

import (
	"context"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Handler interface {
	Kind() domain.JobKind
	Handle(ctx context.Context, job domain.Job) error
}

type Queue interface {
	Claim(ctx context.Context, now time.Time, lease time.Duration) (domain.Job, bool, error)
	Complete(ctx context.Context, job domain.Job) error
	Reschedule(ctx context.Context, job domain.Job, runAt time.Time, cause string) error
	Fail(ctx context.Context, job domain.Job, cause string) error
}

type Clock interface {
	Now() time.Time
}
