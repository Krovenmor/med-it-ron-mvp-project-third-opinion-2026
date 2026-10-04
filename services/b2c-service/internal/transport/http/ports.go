package http

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Graphs interface {
	Build(ctx context.Context, patientID string) (domain.Graph, error)
}

type Bookings interface {
	Slots(ctx context.Context, patientID, recommendationID string) ([]domain.Slot, error)
	Book(ctx context.Context, patientID, recommendationID, slotID string) (domain.Appointment, error)
}
