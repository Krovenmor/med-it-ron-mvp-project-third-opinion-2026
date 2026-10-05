package dashboard

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
)

type Service struct {
	clock Clock
	facts Facts
}

func NewService(clock Clock, facts Facts) *Service {
	return &Service{clock: clock, facts: facts}
}

func (s *Service) Get(ctx context.Context, days int) (domain.Dashboard, error) {
	now := s.clock.Now()
	period, err := domain.NewPeriod(now, days)
	if err != nil {
		return domain.Dashboard{}, err
	}
	from := period.Previous().From

	cases, err := s.facts.CaseFacts(ctx, from, period.To)
	if err != nil {
		return domain.Dashboard{}, err
	}
	tasks, err := s.facts.TaskFacts(ctx, from, period.To)
	if err != nil {
		return domain.Dashboard{}, err
	}
	declines, err := s.facts.DeclineReasons(ctx, period.From, period.To)
	if err != nil {
		return domain.Dashboard{}, err
	}
	return domain.BuildDashboard(period, now, cases, tasks, declines), nil
}
