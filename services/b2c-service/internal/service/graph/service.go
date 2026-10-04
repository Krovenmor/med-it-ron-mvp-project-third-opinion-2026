package graph

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Service struct {
	plans       Plans
	mis         MIS
	clinicPhone string
}

func NewService(plans Plans, mis MIS, clinic config.Clinic) *Service {
	return &Service{plans: plans, mis: mis, clinicPhone: clinic.Phone}
}

func (s *Service) Build(ctx context.Context, patientID string) (domain.Graph, error) {
	plan, err := s.plans.Plan(ctx, patientID)
	if err != nil {
		return domain.Graph{}, err
	}
	history, err := s.mis.History(ctx, patientID)
	if err != nil {
		return domain.Graph{}, err
	}
	return domain.BuildGraph(plan, history, s.clinicPhone), nil
}
