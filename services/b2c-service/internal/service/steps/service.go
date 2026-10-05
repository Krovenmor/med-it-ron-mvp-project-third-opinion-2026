package steps

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Service struct {
	plans Plans
}

func NewService(plans Plans) *Service {
	return &Service{plans: plans}
}

func (s *Service) Decline(ctx context.Context, patientID string, d domain.Decline) error {
	rec, err := s.recommendation(ctx, patientID, d.RecommendationID)
	if err != nil {
		return err
	}
	return s.plans.DeclineRecommendation(ctx, rec.CaseID, d)
}

func (s *Service) RequestHelp(ctx context.Context, patientID, recommendationID string) error {
	rec, err := s.recommendation(ctx, patientID, recommendationID)
	if err != nil {
		return err
	}
	return s.plans.RequestHelp(ctx, rec.CaseID)
}

func (s *Service) recommendation(ctx context.Context, patientID, recommendationID string) (domain.Recommendation, error) {
	plan, err := s.plans.Plan(ctx, patientID)
	if err != nil {
		return domain.Recommendation{}, err
	}
	return plan.Recommendation(recommendationID)
}
