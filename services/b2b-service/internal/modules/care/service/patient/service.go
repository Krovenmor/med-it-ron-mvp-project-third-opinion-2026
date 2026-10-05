package patient

import (
	"context"
	"fmt"
	"slices"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

type Service struct {
	tx        trm.Manager
	clock     Clock
	routes    Routes
	planItems PlanItems
	bookings  Bookings
	tasks     Tasks
	events    Events
	protocol  Protocol
}

func NewService(
	tx trm.Manager,
	clock Clock,
	routes Routes,
	planItems PlanItems,
	bookings Bookings,
	tasks Tasks,
	events Events,
	protocol Protocol,
) *Service {
	return &Service{
		tx:        tx,
		clock:     clock,
		routes:    routes,
		planItems: planItems,
		bookings:  bookings,
		tasks:     tasks,
		events:    events,
		protocol:  protocol,
	}
}

func (s *Service) Decline(ctx context.Context, cmd Decline) (DeclineResult, error) {
	var result DeclineResult
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		route, err := s.routes.GetForUpdate(ctx, cmd.CaseID)
		if err != nil {
			return err
		}
		if err := route.EnsureBookable(); err != nil {
			return err
		}
		item, err := s.planItems.GetForUpdate(ctx, cmd.CaseID, cmd.RecommendationID)
		if err != nil {
			return err
		}
		booked, err := s.bookings.BookedRecommendationIDs(ctx, cmd.CaseID)
		if err != nil {
			return err
		}
		if slices.Contains(booked, item.ID) {
			return fmt.Errorf("%w: recommendation is already booked", apperr.ErrInvalidState)
		}

		now := s.clock.Now()
		if err := item.Decline(cmd.Reason, cmd.Comment, now); err != nil {
			return err
		}
		if err := s.planItems.UpdateDecline(ctx, item); err != nil {
			return err
		}
		if err := s.events.Append(ctx, domain.RecommendationDeclined(item, actorPatient)); err != nil {
			return err
		}
		result = DeclineResult{Item: item, CaseStatus: route.Status}

		open, err := s.hasOpenItems(ctx, cmd.CaseID, booked)
		if err != nil || open {
			return err
		}
		from := route.Status
		if route.Decline(now) {
			if err := s.routes.Update(ctx, route); err != nil {
				return err
			}
			if err := s.events.Append(ctx, domain.StatusChanged(route.CaseID, from, route.Status, actorPatient, now)); err != nil {
				return err
			}
			result.CaseStatus = route.Status
		}
		return s.protocol.AfterDecline(ctx, cmd.CaseID, actorPatient, now)
	})
	return result, err
}

func (s *Service) RequestHelp(ctx context.Context, caseID uuid.UUID) (HelpResult, error) {
	var result HelpResult
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		route, err := s.routes.GetForUpdate(ctx, caseID)
		if err != nil {
			return err
		}
		task, created, err := s.tasks.CreateIfNoActive(ctx, domain.NewHelpRequestTask(route.CaseID, route.Urgency, s.clock.Now()))
		if err != nil {
			return err
		}
		if !created {
			if task, _, err = s.tasks.ActiveByCaseForUpdate(ctx, caseID); err != nil {
				return err
			}
			result = HelpResult{Task: task}
			return nil
		}
		result = HelpResult{Task: task, Created: true}
		return s.events.Append(ctx, domain.OperatorTaskCreated(task))
	})
	return result, err
}

func (s *Service) hasOpenItems(ctx context.Context, caseID uuid.UUID, booked []uuid.UUID) (bool, error) {
	items, err := s.planItems.ListByCase(ctx, caseID)
	if err != nil {
		return false, err
	}
	return slices.ContainsFunc(items, func(i domain.PlanItem) bool {
		return !i.Declined() && !slices.Contains(booked, i.ID)
	}), nil
}
