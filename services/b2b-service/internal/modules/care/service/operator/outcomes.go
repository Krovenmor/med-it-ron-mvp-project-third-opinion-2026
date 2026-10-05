package operator

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/service/booking"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

func (s *Service) Take(ctx context.Context, taskID uuid.UUID, actor string) (domain.OperatorTask, error) {
	return s.withTask(ctx, taskID, func(ctx context.Context, task *domain.OperatorTask, _ *domain.Route, now time.Time) error {
		return task.Take(actor, now)
	})
}

func (s *Service) NoAnswer(ctx context.Context, cmd Outcome) (domain.OperatorTask, error) {
	return s.withTask(ctx, cmd.TaskID, func(ctx context.Context, task *domain.OperatorTask, c *domain.Route, now time.Time) error {
		limitReached, err := task.NoAnswer(cmd.Actor, c.Urgency, now)
		if err != nil {
			return err
		}
		if err := s.record(ctx, c.CaseID, domain.NewCallAttempt(task.ID, domain.CallOutcomeNoAnswer, cmd.Comment, cmd.Actor, now)); err != nil {
			return err
		}
		if !limitReached {
			return nil
		}
		if err := s.notifications.InsertMany(ctx, []domain.Notification{domain.UnreachableNotification(c.CaseID, s.clinicPhone, now)}); err != nil {
			return err
		}
		from := c.Status
		if !c.MarkUnreachable(now) {
			return nil
		}
		return s.saveRoute(ctx, *c, from, cmd.Actor, now)
	})
}

func (s *Service) Callback(ctx context.Context, cmd Callback) (domain.OperatorTask, error) {
	return s.withTask(ctx, cmd.TaskID, func(ctx context.Context, task *domain.OperatorTask, c *domain.Route, now time.Time) error {
		if err := task.Callback(cmd.Actor, cmd.At, now); err != nil {
			return err
		}
		attempt := domain.NewCallAttempt(task.ID, domain.CallOutcomeCallback, cmd.Comment, cmd.Actor, now)
		attempt.CallbackAt = cmd.At
		return s.record(ctx, c.CaseID, attempt)
	})
}

func (s *Service) Contacted(ctx context.Context, cmd Outcome) (domain.OperatorTask, error) {
	return s.withTask(ctx, cmd.TaskID, func(ctx context.Context, task *domain.OperatorTask, c *domain.Route, now time.Time) error {
		attempt := domain.NewCallAttempt(task.ID, domain.CallOutcomeContacted, cmd.Comment, cmd.Actor, now)
		return s.close(ctx, task, c, domain.TaskStatusContacted, attempt, now)
	})
}

func (s *Service) Decline(ctx context.Context, cmd Decline) (domain.OperatorTask, error) {
	return s.withTask(ctx, cmd.TaskID, func(ctx context.Context, task *domain.OperatorTask, c *domain.Route, now time.Time) error {
		attempt, err := domain.NewDeclineAttempt(task.ID, cmd.Reason, cmd.Comment, cmd.Actor, now)
		if err != nil {
			return err
		}
		if err := s.close(ctx, task, c, domain.TaskStatusDeclined, attempt, now); err != nil {
			return err
		}
		from := c.Status
		if c.Decline(now) {
			if err := s.saveRoute(ctx, *c, from, cmd.Actor, now); err != nil {
				return err
			}
		}
		return s.protocol.Stop(ctx, c.CaseID)
	})
}

func (s *Service) HandToDoctor(ctx context.Context, cmd Outcome) (domain.OperatorTask, error) {
	return s.withTask(ctx, cmd.TaskID, func(ctx context.Context, task *domain.OperatorTask, c *domain.Route, now time.Time) error {
		attempt := domain.NewCallAttempt(task.ID, domain.CallOutcomeHandedToDoctor, cmd.Comment, cmd.Actor, now)
		if err := s.close(ctx, task, c, domain.TaskStatusHandedToDoctor, attempt, now); err != nil {
			return err
		}
		patient, err := s.patients.Get(ctx, c.PatientID)
		if err != nil {
			return err
		}
		notification := domain.HandedToDoctorNotification(c.CaseID, patient.ExternalID, cmd.Comment, now)
		return s.notifications.InsertMany(ctx, []domain.Notification{notification})
	})
}

func (s *Service) Book(ctx context.Context, cmd Book) (booking.Result, error) {
	task, err := s.tasks.Get(ctx, cmd.TaskID)
	if err != nil {
		return booking.Result{}, err
	}
	if !task.Status.Active() {
		return booking.Result{}, fmt.Errorf("%w: task is %s", apperr.ErrInvalidState, task.Status)
	}
	route, err := s.routes.Get(ctx, task.CaseID)
	if err != nil {
		return booking.Result{}, err
	}
	if err := route.EnsureBookable(); err != nil {
		return booking.Result{}, err
	}
	item, err := s.planItems.Get(ctx, task.CaseID, cmd.RecommendationID)
	if err != nil {
		return booking.Result{}, err
	}
	if err := item.EnsureBookable(); err != nil {
		return booking.Result{}, err
	}

	appointment, err := s.mis.Book(ctx, gatewayapi.AppointmentRequest{
		PatientID:   route.PatientID,
		SlotID:      cmd.SlotID,
		ServiceName: item.ServiceName,
		ReferralID:  item.ID.String(),
	})
	if err != nil {
		return booking.Result{}, err
	}

	var result booking.Result
	err = s.tx.Do(ctx, func(ctx context.Context) error {
		result, err = s.booker.Book(ctx, booking.Book{
			CaseID:           route.CaseID,
			RecommendationID: item.ID,
			AppointmentID:    appointment.ID,
			ScheduledAt:      appointment.ScheduledAt,
			Channel:          domain.BookingChannelOperator,
		})
		if err != nil {
			return err
		}
		attempt := domain.NewCallAttempt(task.ID, domain.CallOutcomeBooked, cmd.Comment, cmd.Actor, s.clock.Now())
		return s.record(ctx, route.CaseID, attempt)
	})
	return result, err
}

func (s *Service) withTask(
	ctx context.Context,
	taskID uuid.UUID,
	fn func(ctx context.Context, task *domain.OperatorTask, c *domain.Route, now time.Time) error,
) (domain.OperatorTask, error) {
	found, err := s.tasks.Get(ctx, taskID)
	if err != nil {
		return domain.OperatorTask{}, err
	}

	var task domain.OperatorTask
	err = s.tx.Do(ctx, func(ctx context.Context) error {
		c, err := s.routes.GetForUpdate(ctx, found.CaseID)
		if err != nil {
			return err
		}
		task, err = s.tasks.GetForUpdate(ctx, taskID)
		if err != nil {
			return err
		}
		if err := fn(ctx, &task, &c, s.clock.Now()); err != nil {
			return err
		}
		return s.tasks.Update(ctx, task)
	})
	return task, err
}

func (s *Service) close(
	ctx context.Context,
	task *domain.OperatorTask,
	c *domain.Route,
	status domain.TaskStatus,
	attempt domain.CallAttempt,
	now time.Time,
) error {
	if err := task.Close(status, now); err != nil {
		return err
	}
	if err := s.record(ctx, c.CaseID, attempt); err != nil {
		return err
	}
	return s.events.Append(ctx, domain.OperatorTaskClosed(*task, attempt.Actor))
}

func (s *Service) record(ctx context.Context, caseID uuid.UUID, attempt domain.CallAttempt) error {
	if err := s.attempts.Insert(ctx, attempt); err != nil {
		return err
	}
	return s.events.Append(ctx, domain.CallAttempted(caseID, attempt))
}

func (s *Service) saveRoute(ctx context.Context, c domain.Route, from domain.RouteStatus, actor string, now time.Time) error {
	if err := s.routes.Update(ctx, c); err != nil {
		return err
	}
	return s.events.Append(ctx, domain.StatusChanged(c.CaseID, from, c.Status, actor, now))
}
