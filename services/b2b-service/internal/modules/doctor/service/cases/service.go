package cases

import (
	"context"

	"github.com/google/uuid"

	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
)

type Service struct {
	cases           Cases
	recommendations Recommendations
	mis             MIS
}

func NewService(cases Cases, recommendations Recommendations, mis MIS) *Service {
	return &Service{cases: cases, recommendations: recommendations, mis: mis}
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (Details, error) {
	c, err := s.cases.Get(ctx, id)
	if err != nil {
		return Details{}, err
	}
	recs, err := s.recommendations.ListByCase(ctx, id)
	if err != nil {
		return Details{}, err
	}
	return Details{Case: c, Recommendations: recs}, nil
}

func (s *Service) History(ctx context.Context, id uuid.UUID) (gatewayapi.History, error) {
	c, err := s.cases.Get(ctx, id)
	if err != nil {
		return gatewayapi.History{}, err
	}
	return s.mis.PatientHistory(ctx, c.Patient.ID)
}
