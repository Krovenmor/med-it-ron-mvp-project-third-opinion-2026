package booking

import (
	"context"
	"fmt"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Service struct {
	plans Plans
	mis   MIS
}

func NewService(plans Plans, mis MIS) *Service {
	return &Service{plans: plans, mis: mis}
}

func (s *Service) Slots(ctx context.Context, patientID, recommendationID string) ([]domain.Slot, error) {
	rec, err := s.bookableRecommendation(ctx, patientID, recommendationID)
	if err != nil {
		return nil, err
	}
	return s.mis.Slots(ctx, rec.ServiceCode)
}

func (s *Service) Book(ctx context.Context, patientID, recommendationID, slotID string) (domain.Appointment, error) {
	if slotID == "" {
		return domain.Appointment{}, fmt.Errorf("%w: slot_id is required", domain.ErrInvalidInput)
	}
	rec, err := s.bookableRecommendation(ctx, patientID, recommendationID)
	if err != nil {
		return domain.Appointment{}, err
	}

	history, err := s.mis.History(ctx, patientID)
	if err != nil {
		return domain.Appointment{}, err
	}
	if existing, booked := history.AppointmentByReferral(rec.ID); booked {
		if err := s.register(ctx, rec, existing); err != nil {
			return domain.Appointment{}, err
		}
		return domain.Appointment{}, fmt.Errorf("%w: recommendation is already booked", domain.ErrInvalidState)
	}

	appointment, err := s.mis.Book(ctx, domain.AppointmentRequest{
		PatientID:   patientID,
		SlotID:      slotID,
		ServiceName: rec.ServiceName,
		ReferralID:  rec.ID,
	})
	if err != nil {
		return domain.Appointment{}, err
	}
	if err := s.register(ctx, rec, appointment); err != nil {
		return domain.Appointment{}, err
	}
	return appointment, nil
}

func (s *Service) bookableRecommendation(ctx context.Context, patientID, recommendationID string) (domain.Recommendation, error) {
	plan, err := s.plans.Plan(ctx, patientID)
	if err != nil {
		return domain.Recommendation{}, err
	}
	rec, err := plan.Recommendation(recommendationID)
	if err != nil {
		return domain.Recommendation{}, err
	}
	return rec, rec.EnsureBookable()
}

func (s *Service) register(ctx context.Context, rec domain.Recommendation, appointment domain.Appointment) error {
	return s.plans.RegisterBooking(ctx, rec.CaseID, domain.Booking{
		RecommendationID: rec.ID,
		AppointmentID:    appointment.ID,
		ScheduledAt:      appointment.ScheduledAt,
	})
}
