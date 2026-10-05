package routes

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
)

type Clock interface {
	Now() time.Time
}

type Patients interface {
	Upsert(ctx context.Context, p domain.Patient, now time.Time) error
}

type Routes interface {
	CreateIfAbsent(ctx context.Context, r domain.Route) (bool, error)
	GetForUpdate(ctx context.Context, caseID uuid.UUID) (domain.Route, error)
	Update(ctx context.Context, r domain.Route) error
}

type PlanItems interface {
	InsertMany(ctx context.Context, items []domain.PlanItem) error
	ActiveServices(ctx context.Context, patientID, excludeCaseID uuid.UUID) ([]domain.Service, error)
}

type Notifications interface {
	InsertMany(ctx context.Context, notifications []domain.Notification) error
}

type Events interface {
	Append(ctx context.Context, e domain.CaseEvent) error
}

type Protocol interface {
	StartEmergency(ctx context.Context, caseID uuid.UUID, now time.Time) error
	StopEmergency(ctx context.Context, caseID uuid.UUID, actor string, now time.Time) error
	AfterConfirm(ctx context.Context, r domain.Route, now time.Time) error
}
