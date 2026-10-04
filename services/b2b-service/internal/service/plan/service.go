package plan

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Service struct {
	patients        Patients
	cases           Cases
	recommendations Recommendations
}

func NewService(patients Patients, cases Cases, recommendations Recommendations) *Service {
	return &Service{patients: patients, cases: cases, recommendations: recommendations}
}

func (s *Service) Get(ctx context.Context, sourceSystem, externalID string) (domain.PatientPlan, error) {
	patient, err := s.patients.GetByExternalID(ctx, sourceSystem, externalID)
	if err != nil {
		return domain.PatientPlan{}, err
	}
	cases, err := s.cases.ListByPatient(ctx, patient.ID)
	if err != nil {
		return domain.PatientPlan{}, err
	}

	plan := domain.PatientPlan{Patient: patient, Cases: make([]domain.PlanCase, 0, len(cases))}
	for _, c := range cases {
		if !c.Confirmed() && !c.NeedsUrgentContact() {
			continue
		}
		planCase := domain.PlanCase{Case: c, UrgentContact: c.NeedsUrgentContact()}
		if c.Confirmed() {
			if planCase.Recommendations, err = s.visibleRecommendations(ctx, c); err != nil {
				return domain.PatientPlan{}, err
			}
		}
		plan.Cases = append(plan.Cases, planCase)
	}
	return plan, nil
}

func (s *Service) visibleRecommendations(ctx context.Context, c domain.Case) ([]domain.Recommendation, error) {
	recs, err := s.recommendations.ListByCase(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	visible := make([]domain.Recommendation, 0, len(recs))
	for _, rec := range recs {
		if rec.VisibleToPatient() {
			visible = append(visible, rec)
		}
	}
	return visible, nil
}
