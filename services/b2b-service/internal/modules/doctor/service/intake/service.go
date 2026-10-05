package intake

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/events"
)

type Service struct {
	cases           Cases
	recommendations Recommendations
	events          Events
	jobs            Jobs
	publisher       Publisher
}

func NewService(cases Cases, recommendations Recommendations, events Events, jobs Jobs, publisher Publisher) *Service {
	return &Service{cases: cases, recommendations: recommendations, events: events, jobs: jobs, publisher: publisher}
}

func (s *Service) Subscription(consumer *events.Consumer) events.Subscription {
	return consumer.Subscribe(gatewayapi.TopicCaseAssessed, func(ctx context.Context, m events.Message) error {
		var e gatewayapi.CaseAssessed
		if err := m.Decode(&e); err != nil {
			return err
		}
		return s.Accept(ctx, e)
	})
}

func (s *Service) Accept(ctx context.Context, e gatewayapi.CaseAssessed) error {
	c := caseOf(e)
	created, err := s.cases.CreateIfAbsent(ctx, c)
	if err != nil || !created {
		return err
	}
	if err := s.recommendations.InsertMany(ctx, recommendationsOf(c.ID, e)); err != nil {
		return err
	}
	if err := s.events.Append(ctx, domain.StatusChanged(c.ID, "", c.Status, domain.ActorSystem, c.AssessedAt)); err != nil {
		return err
	}
	if c.Urgency == domain.UrgencyEmergency {
		if err := s.jobs.Enqueue(ctx, domain.ReviewEscalationJob(c.ID, c.AssessedAt), c.AssessedAt); err != nil {
			return err
		}
	}
	return s.publisher.Publish(ctx, api.TopicReviewStarted, c.ID.String(), reviewStarted(c))
}
