package booking

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
)

type Clock interface {
	Now() time.Time
}

type Routes interface {
	GetForUpdate(ctx context.Context, caseID uuid.UUID) (domain.Route, error)
	Update(ctx context.Context, r domain.Route) error
}

type PlanItems interface {
	Get(ctx context.Context, caseID, id uuid.UUID) (domain.PlanItem, error)
}

type Bookings interface {
	CreateIfAbsent(ctx context.Context, b domain.Booking) (domain.Booking, bool, error)
}

type Events interface {
	Append(ctx context.Context, e domain.CaseEvent) error
}

type Protocol interface {
	AfterBooking(ctx context.Context, caseID uuid.UUID, actor string, now time.Time) error
}
