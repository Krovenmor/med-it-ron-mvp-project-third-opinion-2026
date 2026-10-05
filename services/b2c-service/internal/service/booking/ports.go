package booking

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Plans interface {
	Plan(ctx context.Context, patientID string) (domain.Plan, error)
	RegisterBooking(ctx context.Context, caseID string, b domain.Booking) error
}

type MIS interface {
	History(ctx context.Context, patientID string) (domain.History, error)
	Slots(ctx context.Context, serviceCode string) ([]domain.Slot, error)
	Book(ctx context.Context, req domain.AppointmentRequest) (domain.Appointment, error)
}
