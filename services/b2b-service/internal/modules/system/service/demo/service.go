package demo

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
)

type Service struct {
	clock   Clock
	modules []Module
}

func NewService(clock Clock, modules ...Module) *Service {
	return &Service{clock: clock, modules: modules}
}

func (s *Service) Reset(ctx context.Context) error {
	s.clock.Reset()
	now := s.clock.Now()
	cases := append(history.Generate(now), history.DemoScenario(now)...)
	for _, m := range s.modules {
		if err := m.Reset(ctx, cases); err != nil {
			return err
		}
	}
	return nil
}
