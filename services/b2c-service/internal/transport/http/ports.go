package http

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Routes interface {
	Build(ctx context.Context, patientID string) (domain.Route, error)
}

type Bookings interface {
	Slots(ctx context.Context, patientID, recommendationID string) ([]domain.Slot, error)
	Book(ctx context.Context, patientID, recommendationID, slotID string) (domain.Appointment, error)
}

type Steps interface {
	Decline(ctx context.Context, patientID string, d domain.Decline) error
	RequestHelp(ctx context.Context, patientID, recommendationID string) error
}
