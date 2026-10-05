package mis

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
)

type Patients interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Patient, error)
}

type Client interface {
	PatientHistory(ctx context.Context, externalID string) (domain.PatientHistory, error)
	Slots(ctx context.Context, serviceCode string, from time.Time) ([]domain.Slot, error)
	Book(ctx context.Context, req domain.AppointmentRequest) (domain.Appointment, error)
}
