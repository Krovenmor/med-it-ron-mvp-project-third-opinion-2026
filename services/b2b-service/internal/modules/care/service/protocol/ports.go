package protocol

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

type Clock interface {
	Now() time.Time
}

type Routes interface {
	GetForUpdate(ctx context.Context, caseID uuid.UUID) (domain.Route, error)
	Update(ctx context.Context, r domain.Route) error
}

type Patients interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Patient, error)
}

type Tasks interface {
	CreateIfNoActive(ctx context.Context, t domain.OperatorTask) (domain.OperatorTask, bool, error)
	ActiveByCaseForUpdate(ctx context.Context, caseID uuid.UUID) (domain.OperatorTask, bool, error)
	Update(ctx context.Context, t domain.OperatorTask) error
}

type Jobs interface {
	Enqueue(ctx context.Context, job jobs.Job, now time.Time) error
	CancelPending(ctx context.Context, key string, kinds []string) error
}

type Notifications interface {
	InsertMany(ctx context.Context, notifications []domain.Notification) error
}

type Events interface {
	Append(ctx context.Context, e domain.CaseEvent) error
}
