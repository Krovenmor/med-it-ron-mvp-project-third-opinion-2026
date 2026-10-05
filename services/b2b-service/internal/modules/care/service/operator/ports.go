package operator

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/booking"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
)

type Clock interface {
	Now() time.Time
}

type Tasks interface {
	Get(ctx context.Context, id uuid.UUID) (domain.OperatorTask, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (domain.OperatorTask, error)
	Update(ctx context.Context, t domain.OperatorTask) error
	ListActive(ctx context.Context, now time.Time, assignee string) ([]domain.TaskListItem, error)
}

type Attempts interface {
	Insert(ctx context.Context, a domain.CallAttempt) error
	ListByTask(ctx context.Context, taskID uuid.UUID) ([]domain.CallAttempt, error)
}

type Routes interface {
	Get(ctx context.Context, caseID uuid.UUID) (domain.Route, error)
	GetForUpdate(ctx context.Context, caseID uuid.UUID) (domain.Route, error)
	Update(ctx context.Context, r domain.Route) error
}

type Patients interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Patient, error)
}

type PlanItems interface {
	Get(ctx context.Context, caseID, id uuid.UUID) (domain.PlanItem, error)
	ListByCase(ctx context.Context, caseID uuid.UUID) ([]domain.PlanItem, error)
}

type Bookings interface {
	BookedRecommendationIDs(ctx context.Context, caseID uuid.UUID) ([]uuid.UUID, error)
}

type Notifications interface {
	InsertMany(ctx context.Context, notifications []domain.Notification) error
	List(ctx context.Context, limit int) ([]domain.NotificationEntry, error)
}

type Events interface {
	Append(ctx context.Context, e domain.CaseEvent) error
}

type Protocol interface {
	Stop(ctx context.Context, caseID uuid.UUID) error
}

type MIS interface {
	Slots(ctx context.Context, serviceCode string, from time.Time) ([]gatewayapi.Slot, error)
	Book(ctx context.Context, req gatewayapi.AppointmentRequest) (gatewayapi.Appointment, error)
}

type Booker interface {
	Book(ctx context.Context, cmd booking.Book) (booking.Result, error)
}
