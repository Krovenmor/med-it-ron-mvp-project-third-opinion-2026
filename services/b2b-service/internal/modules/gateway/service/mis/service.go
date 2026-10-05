package mis

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
)

type Service struct {
	patients Patients
	client   Client
}

func NewService(patients Patients, client Client) *Service {
	return &Service{patients: patients, client: client}
}

func (s *Service) PatientHistory(ctx context.Context, patientID uuid.UUID) (api.History, error) {
	patient, err := s.patients.Get(ctx, patientID)
	if err != nil {
		return api.History{}, err
	}
	h, err := s.client.PatientHistory(ctx, patient.ExternalID)
	if err != nil {
		return api.History{}, err
	}
	return historyOf(h), nil
}

func (s *Service) Slots(ctx context.Context, serviceCode string, from time.Time) ([]api.Slot, error) {
	slots, err := s.client.Slots(ctx, serviceCode, from)
	if err != nil {
		return nil, err
	}
	return slotsOf(slots), nil
}

func (s *Service) Book(ctx context.Context, req api.AppointmentRequest) (api.Appointment, error) {
	patient, err := s.patients.Get(ctx, req.PatientID)
	if err != nil {
		return api.Appointment{}, err
	}
	appointment, err := s.client.Book(ctx, domain.AppointmentRequest{
		PatientID:   patient.ExternalID,
		SlotID:      req.SlotID,
		ServiceName: req.ServiceName,
		ReferralID:  req.ReferralID,
	})
	if err != nil {
		return api.Appointment{}, err
	}
	return appointmentOf(appointment), nil
}
