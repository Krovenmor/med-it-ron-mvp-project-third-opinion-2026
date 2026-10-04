package cases

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	cases           Cases
	patients        Patients
	recommendations Recommendations
}

func NewService(cases Cases, patients Patients, recommendations Recommendations) *Service {
	return &Service{cases: cases, patients: patients, recommendations: recommendations}
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
