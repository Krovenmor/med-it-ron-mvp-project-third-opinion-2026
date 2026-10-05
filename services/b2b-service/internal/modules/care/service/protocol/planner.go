package protocol

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

type Planner struct {
	tasks         Tasks
	jobs          Jobs
	notifications Notifications
	events        Events
	clinicPhone   string
}

func NewPlanner(tasks Tasks, jobs Jobs, notifications Notifications, events Events, clinic config.Clinic) *Planner {
	return &Planner{tasks: tasks, jobs: jobs, notifications: notifications, events: events, clinicPhone: clinic.Phone}
}

func (p *Planner) StartEmergency(ctx context.Context, caseID uuid.UUID, now time.Time) error {
	task, created, err := p.tasks.CreateIfNoActive(ctx, domain.NewEmergencyTask(caseID, now))
	if err != nil || !created {
		return err
	}
	if err := p.events.Append(ctx, domain.OperatorTaskCreated(task)); err != nil {
		return err
	}
	if err := p.notifications.InsertMany(ctx, domain.UrgentContactNotifications(caseID, p.clinicPhone, now)); err != nil {
		return err
	}
	return p.enqueue(ctx, domain.EmergencyJobs(caseID, now), now)
}

func (p *Planner) StopEmergency(ctx context.Context, caseID uuid.UUID, actor string, now time.Time) error {
	if err := p.jobs.CancelPending(ctx, caseID.String(), domain.EmergencyJobKinds); err != nil {
		return err
	}
	return p.closeActiveTask(ctx, caseID, domain.TaskReasonEmergency, domain.TaskStatusCancelled, actor, now)
}

func (p *Planner) AfterConfirm(ctx context.Context, r domain.Route, now time.Time) error {
	return p.enqueue(ctx, domain.JobsAfterConfirm(r, now), now)
}

func (p *Planner) AfterBooking(ctx context.Context, caseID uuid.UUID, actor string, now time.Time) error {
	if err := p.Stop(ctx, caseID); err != nil {
		return err
	}
	return p.closeActiveTask(ctx, caseID, "", domain.TaskStatusBooked, actor, now)
}

func (p *Planner) AfterDecline(ctx context.Context, caseID uuid.UUID, actor string, now time.Time) error {
	if err := p.Stop(ctx, caseID); err != nil {
		return err
	}
	return p.closeActiveTask(ctx, caseID, "", domain.TaskStatusDeclined, actor, now)
}

func (p *Planner) Stop(ctx context.Context, caseID uuid.UUID) error {
	return p.jobs.CancelPending(ctx, caseID.String(), domain.ProtocolJobKinds)
}

func (p *Planner) closeActiveTask(
	ctx context.Context,
	caseID uuid.UUID,
	reason domain.TaskReason,
	status domain.TaskStatus,
	actor string,
	now time.Time,
) error {
	task, found, err := p.tasks.ActiveByCaseForUpdate(ctx, caseID)
	if err != nil || !found || (reason != "" && task.Reason != reason) {
		return err
	}
	if err := task.Close(status, now); err != nil {
		return err
	}
	if err := p.tasks.Update(ctx, task); err != nil {
		return err
	}
	return p.events.Append(ctx, domain.OperatorTaskClosed(task, actor))
}

func (p *Planner) enqueue(ctx context.Context, planned []jobs.Job, now time.Time) error {
	for _, job := range planned {
		if err := p.jobs.Enqueue(ctx, job, now); err != nil {
			return err
		}
	}
	return nil
}
