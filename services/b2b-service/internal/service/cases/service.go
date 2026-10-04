package cases

import (
	"context"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Service struct {
	cases           Cases
	patients        Patients
	recommendations Recommendations
	mis             MIS
}

func NewService(cases Cases, patients Patients, recommendations Recommendations, mis MIS) *Service {
	return &Service{cases: cases, patients: patients, recommendations: recommendations, mis: mis}
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (Details, error) {
	c, err := s.cases.Get(ctx, id)
	if err != nil {
		return Details{}, err
	}
	patient, err := s.patients.Get(ctx, c.PatientID)
	if err != nil {
		return Details{}, err
	}
	recs, err := s.recommendations.ListByCase(ctx, id)
	if err != nil {
		return Details{}, err
	}
	return Details{Case: c, Patient: patient, Recommendations: recs}, nil
}

func (s *Service) History(ctx context.Context, id uuid.UUID) (domain.PatientHistory, error) {
	c, err := s.cases.Get(ctx, id)
	if err != nil {
		return domain.PatientHistory{}, err
	}
	patient, err := s.patients.Get(ctx, c.PatientID)
	if err != nil {
		return domain.PatientHistory{}, err
	}
	return s.mis.PatientHistory(ctx, patient)
}
