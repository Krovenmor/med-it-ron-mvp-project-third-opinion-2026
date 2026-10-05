package booking

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

type Service struct {
	tx        trm.Manager
	clock     Clock
	routes    Routes
	planItems PlanItems
	bookings  Bookings
	events    Events
	protocol  Protocol
}

func NewService(
	tx trm.Manager,
	clock Clock,
	routes Routes,
	planItems PlanItems,
	bookings Bookings,
	events Events,
	protocol Protocol,
) *Service {
	return &Service{
		tx:        tx,
		clock:     clock,
		routes:    routes,
		planItems: planItems,
		bookings:  bookings,
		events:    events,
		protocol:  protocol,
	}
}

func (s *Service) Book(ctx context.Context, cmd Book) (Result, error) {
	booking := domain.Booking{
		CaseID:           cmd.CaseID,
		RecommendationID: cmd.RecommendationID,
		AppointmentID:    cmd.AppointmentID,
		ScheduledAt:      cmd.ScheduledAt,
		Channel:          cmd.Channel,
		CreatedAt:        s.clock.Now(),
	}
	if err := booking.Validate(); err != nil {
		return Result{}, err
	}

	var result Result
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		c, err := s.routes.GetForUpdate(ctx, cmd.CaseID)
		if err != nil {
			return err
		}
		if err := c.EnsureBookable(); err != nil {
			return err
		}
		item, err := s.planItems.Get(ctx, cmd.CaseID, cmd.RecommendationID)
		if err != nil {
			return err
		}
		if err := item.EnsureBookable(); err != nil {
			return err
		}

		from := c.Status
		statusChanged, err := c.Book(booking.CreatedAt)
		if err != nil {
			return err
		}

		saved, created, err := s.bookings.CreateIfAbsent(ctx, booking)
		if err != nil {
			return err
		}
		if !created {
			if !saved.SameTarget(booking) {
				return fmt.Errorf("%w: appointment %s is already linked to another recommendation", apperr.ErrConflict, booking.AppointmentID)
			}
			result = Result{Booking: saved, CaseStatus: from}
			return nil
		}

		if statusChanged {
			if err := s.routes.Update(ctx, c); err != nil {
				return err
			}
			if err := s.events.Append(ctx, domain.StatusChanged(c.CaseID, from, c.Status, booking.Channel.Actor(), booking.CreatedAt)); err != nil {
				return err
			}
		}
		if err := s.events.Append(ctx, domain.RecommendationBooked(saved)); err != nil {
			return err
		}
		if err := s.protocol.AfterBooking(ctx, c.CaseID, booking.Channel.Actor(), booking.CreatedAt); err != nil {
			return err
		}
		result = Result{Booking: saved, CaseStatus: c.Status, Created: true}
		return nil
	})
	return result, err
}
