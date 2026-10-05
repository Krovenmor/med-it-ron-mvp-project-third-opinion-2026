package routes

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	careapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	doctorapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/api"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/events"
)

var errReviewNotStarted = fmt.Errorf("%w: review has not started yet", apperr.ErrInvalidState)

type Service struct {
	clock         Clock
	patients      Patients
	routes        Routes
	planItems     PlanItems
	notifications Notifications
	events        Events
	protocol      Protocol
}

func NewService(
	clock Clock,
	patients Patients,
	routes Routes,
	planItems PlanItems,
	notifications Notifications,
	events Events,
	protocol Protocol,
) *Service {
	return &Service{
		clock:         clock,
		patients:      patients,
		routes:        routes,
		planItems:     planItems,
		notifications: notifications,
		events:        events,
		protocol:      protocol,
	}
}

func (s *Service) Subscriptions(consumer *events.Consumer) []events.Subscription {
	return []events.Subscription{
		consumer.Subscribe(gatewayapi.TopicPatientRegistered, handle(s.PatientRegistered)),
		consumer.Subscribe(gatewayapi.TopicReportReceived, handle(s.ReportReceived)),
		consumer.Subscribe(doctorapi.TopicReviewStarted, handle(s.ReviewStarted)),
		consumer.Subscribe(doctorapi.TopicUrgencyChanged, handle(s.UrgencyChanged)),
		consumer.Subscribe(doctorapi.TopicCaseConfirmed, handle(s.CaseConfirmed)),
		consumer.Subscribe(doctorapi.TopicReviewEscalated, handle(s.ReviewEscalated)),
	}
}

func (s *Service) ActiveServices(ctx context.Context, patientID, excludeCaseID uuid.UUID) ([]careapi.Service, error) {
	services, err := s.planItems.ActiveServices(ctx, patientID, excludeCaseID)
	if err != nil {
		return nil, err
	}
	return servicesOf(services), nil
}

func (s *Service) PatientRegistered(ctx context.Context, e gatewayapi.PatientRegistered) error {
	return s.patients.Upsert(ctx, patientOf(e), s.clock.Now())
}

func (s *Service) ReportReceived(ctx context.Context, e gatewayapi.ReportReceived) error {
	route := domain.NewRoute(e.CaseID, e.PatientID, domain.Modality(e.Modality), e.PerformedAt, e.ReceivedAt)
	created, err := s.routes.CreateIfAbsent(ctx, route)
	if err != nil || !created {
		return err
	}
	return s.events.Append(ctx, domain.StatusChanged(route.CaseID, "", route.Status, domain.ActorSystem, route.ReceivedAt))
}

func (s *Service) ReviewStarted(ctx context.Context, e doctorapi.ReviewStarted) error {
	start := reviewStartOf(e)
	route, err := s.routes.GetForUpdate(ctx, e.CaseID)
	switch {
	case errors.Is(err, apperr.ErrNotFound):
		route = domain.NewRoute(e.CaseID, e.PatientID, start.Modality, e.PerformedAt, e.ReceivedAt)
		if _, err := s.routes.CreateIfAbsent(ctx, route); err != nil {
			return err
		}
	case err != nil:
		return err
	}

	from := route.Status
	if !route.StartReview(start) {
		return nil
	}
	if err := s.routes.Update(ctx, route); err != nil {
		return err
	}
	if err := s.events.Append(ctx, domain.StatusChanged(route.CaseID, from, route.Status, domain.ActorSystem, e.AssessedAt)); err != nil {
		return err
	}
	if route.Urgency == domain.UrgencyEmergency {
		return s.protocol.StartEmergency(ctx, route.CaseID, e.AssessedAt)
	}
	return nil
}

func (s *Service) UrgencyChanged(ctx context.Context, e doctorapi.UrgencyChanged) error {
	route, err := s.routes.GetForUpdate(ctx, e.CaseID)
	if err != nil {
		return err
	}
	if route.Status == domain.RouteStatusReceived {
		return errReviewNotStarted
	}
	from, to := domain.Urgency(e.From), domain.Urgency(e.To)
	route.ChangeUrgency(to, e.ReviewDueAt, e.ChangedAt)
	if err := s.routes.Update(ctx, route); err != nil {
		return err
	}
	switch {
	case to == domain.UrgencyEmergency:
		return s.protocol.StartEmergency(ctx, route.CaseID, e.ChangedAt)
	case from == domain.UrgencyEmergency:
		return s.protocol.StopEmergency(ctx, route.CaseID, e.Actor, e.ChangedAt)
	}
	return nil
}

func (s *Service) CaseConfirmed(ctx context.Context, e doctorapi.CaseConfirmed) error {
	route, err := s.routes.GetForUpdate(ctx, e.CaseID)
	if err != nil {
		return err
	}
	if route.Status == domain.RouteStatusReceived {
		return errReviewNotStarted
	}
	from := route.Status
	if !route.Confirm(domain.Urgency(e.Urgency), e.ConfirmedAt) {
		return nil
	}
	if err := s.planItems.InsertMany(ctx, planItemsOf(e)); err != nil {
		return err
	}
	if err := s.routes.Update(ctx, route); err != nil {
		return err
	}
	if err := s.events.Append(ctx, domain.StatusChanged(route.CaseID, from, route.Status, e.Actor, e.ConfirmedAt)); err != nil {
		return err
	}
	return s.protocol.AfterConfirm(ctx, route, e.ConfirmedAt)
}

func (s *Service) ReviewEscalated(ctx context.Context, e doctorapi.ReviewEscalated) error {
	notification := domain.ReviewEscalationNotification(e.CaseID, e.PatientExternalID, e.EscalatedAt)
	return s.notifications.InsertMany(ctx, []domain.Notification{notification})
}

func handle[T any](fn func(ctx context.Context, e T) error) events.Handler {
	return func(ctx context.Context, m events.Message) error {
		var e T
		if err := m.Decode(&e); err != nil {
			return err
		}
		return fn(ctx, e)
	}
}
