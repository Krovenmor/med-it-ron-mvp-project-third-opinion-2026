package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/booking"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/operator"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/patient"
)

type Plan interface {
	Get(ctx context.Context, sourceSystem, externalID string) (domain.PatientPlan, error)
}

type Bookings interface {
	Book(ctx context.Context, cmd booking.Book) (booking.Result, error)
}

type Patient interface {
	Decline(ctx context.Context, cmd patient.Decline) (patient.DeclineResult, error)
	RequestHelp(ctx context.Context, caseID uuid.UUID) (patient.HelpResult, error)
}

type Operator interface {
	List(ctx context.Context, assignee string) ([]domain.TaskListItem, error)
	Card(ctx context.Context, taskID uuid.UUID) (operator.Card, error)
	Take(ctx context.Context, taskID uuid.UUID, actor string) (domain.OperatorTask, error)
	NoAnswer(ctx context.Context, cmd operator.Outcome) (domain.OperatorTask, error)
	Callback(ctx context.Context, cmd operator.Callback) (domain.OperatorTask, error)
	Contacted(ctx context.Context, cmd operator.Outcome) (domain.OperatorTask, error)
	Decline(ctx context.Context, cmd operator.Decline) (domain.OperatorTask, error)
	HandToDoctor(ctx context.Context, cmd operator.Outcome) (domain.OperatorTask, error)
	Book(ctx context.Context, cmd operator.Book) (booking.Result, error)
	Journal(ctx context.Context) ([]domain.NotificationEntry, error)
}

type Dashboard interface {
	Get(ctx context.Context, days int) (domain.Dashboard, error)
}
