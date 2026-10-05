package protocol

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

type Handlers struct {
	tx            trm.Manager
	clock         Clock
	routes        Routes
	patients      Patients
	tasks         Tasks
	notifications Notifications
	events        Events
	clinicPhone   string
}

func NewHandlers(
	tx trm.Manager,
	clock Clock,
	routes Routes,
	patients Patients,
	tasks Tasks,
	notifications Notifications,
	events Events,
	clinic config.Clinic,
) *Handlers {
	return &Handlers{
		tx:            tx,
		clock:         clock,
		routes:        routes,
		patients:      patients,
		tasks:         tasks,
		notifications: notifications,
		events:        events,
		clinicPhone:   clinic.Phone,
	}
}

func (h *Handlers) All() []jobs.Handler {
	return []jobs.Handler{
		jobs.Handle(domain.JobKindNotifyPatient, h.withRoute(h.notifyPatient)),
		jobs.Handle(domain.JobKindRemindPatient, h.withRoute(h.remindPatient)),
		jobs.Handle(domain.JobKindFollowUp, h.withRoute(h.followUp)),
		jobs.Handle(domain.JobKindEscalateContact, h.withRoute(h.escalateContact)),
	}
}

func (h *Handlers) notifyPatient(ctx context.Context, r domain.Route) error {
	now := h.clock.Now()
	if r.Status != domain.RouteStatusConfirmed {
		return nil
	}
	if err := h.notifications.InsertMany(ctx, domain.PlanReadyNotifications(r, h.clinicPhone, now)); err != nil {
		return err
	}
	from := r.Status
	r.MarkNotified(now)
	if err := h.routes.Update(ctx, r); err != nil {
		return err
	}
	return h.events.Append(ctx, domain.StatusChanged(r.CaseID, from, r.Status, domain.ActorSystem, now))
}

func (h *Handlers) remindPatient(ctx context.Context, r domain.Route) error {
	if !r.AwaitingBooking() {
		return nil
	}
	return h.notifications.InsertMany(ctx, []domain.Notification{domain.ReminderNotification(r.CaseID, h.clock.Now())})
}

func (h *Handlers) followUp(ctx context.Context, r domain.Route) error {
	if !r.AwaitingBooking() {
		return nil
	}
	task, created, err := h.tasks.CreateIfNoActive(ctx, domain.NewFollowUpTask(r.CaseID, r.Urgency, h.clock.Now()))
	if err != nil || !created {
		return err
	}
	return h.events.Append(ctx, domain.OperatorTaskCreated(task))
}

func (h *Handlers) escalateContact(ctx context.Context, r domain.Route) error {
	task, found, err := h.tasks.ActiveByCaseForUpdate(ctx, r.CaseID)
	if err != nil || !found || task.Reason != domain.TaskReasonEmergency || task.Status == domain.TaskStatusCallback {
		return err
	}
	patient, err := h.patients.Get(ctx, r.PatientID)
	if err != nil {
		return err
	}
	notification := domain.ContactEscalationNotification(r.CaseID, patient.ExternalID, h.clock.Now())
	return h.notifications.InsertMany(ctx, []domain.Notification{notification})
}

func (h *Handlers) withRoute(fn func(ctx context.Context, r domain.Route) error) func(ctx context.Context, job jobs.Job) error {
	return func(ctx context.Context, job jobs.Job) error {
		caseID, err := uuid.Parse(job.Key)
		if err != nil {
			return fmt.Errorf("%s job key: %w", job.Kind, err)
		}
		return h.tx.Do(ctx, func(ctx context.Context) error {
			r, err := h.routes.GetForUpdate(ctx, caseID)
			if err != nil {
				return err
			}
			return fn(ctx, r)
		})
	}
}
