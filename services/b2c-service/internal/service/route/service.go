package route

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Service struct {
	plans  Plans
	mis    MIS
	clinic domain.Clinic
}

func NewService(plans Plans, mis MIS, clinic config.Clinic) *Service {
	return &Service{plans: plans, mis: mis, clinic: domain.Clinic{Name: clinic.Name, Phone: clinic.Phone, Address: clinic.Address}}
}

func (s *Service) Build(ctx context.Context, patientID string) (domain.Route, error) {
	plan, err := s.plans.Plan(ctx, patientID)
	if err != nil {
		return domain.Route{}, err
	}
	history, err := s.mis.History(ctx, patientID)
	if err != nil {
		return domain.Route{}, err
	}
	return domain.BuildRoute(plan, history, s.clinic), nil
}
