package plan

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
)

type Service struct {
	clock     Clock
	patients  Patients
	routes    Routes
	planItems PlanItems
}

func NewService(clock Clock, patients Patients, routes Routes, planItems PlanItems) *Service {
	return &Service{clock: clock, patients: patients, routes: routes, planItems: planItems}
}

func (s *Service) Get(ctx context.Context, sourceSystem, externalID string) (domain.PatientPlan, error) {
	patient, err := s.patients.GetByExternalID(ctx, sourceSystem, externalID)
	if err != nil {
		return domain.PatientPlan{}, err
	}
	routes, err := s.routes.ListByPatient(ctx, patient.ID)
	if err != nil {
		return domain.PatientPlan{}, err
	}

	plan := domain.PatientPlan{Now: s.clock.Now(), Patient: patient, Cases: make([]domain.PlanCase, 0, len(routes))}
	for _, r := range routes {
		if !r.Confirmed() && !r.InReview() {
			continue
		}
		planCase := domain.PlanCase{Route: r, UrgentContact: r.NeedsUrgentContact()}
		if r.Confirmed() {
			if planCase.Items, err = s.planItems.ListByCase(ctx, r.CaseID); err != nil {
				return domain.PatientPlan{}, err
			}
		}
		plan.Cases = append(plan.Cases, planCase)
	}
	return plan, nil
}
