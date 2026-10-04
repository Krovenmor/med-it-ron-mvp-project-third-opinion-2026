package booking

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
	GetForUpdate(ctx context.Context, id uuid.UUID) (domain.Case, error)
	Update(ctx context.Context, c domain.Case) error
}

type Recommendations interface {
	Get(ctx context.Context, caseID, id uuid.UUID) (domain.Recommendation, error)
}

type Bookings interface {
	CreateIfAbsent(ctx context.Context, b domain.Booking) (domain.Booking, bool, error)
}

type Events interface {
	Append(ctx context.Context, e domain.CaseEvent) error
}
