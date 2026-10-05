package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

type TaskReason string

const (
	TaskReasonEmergency   TaskReason = "emergency"
	TaskReasonNoBooking   TaskReason = "no_booking"
	TaskReasonHelpRequest TaskReason = "help_request"
)

type TaskStatus string

const (
	TaskStatusNew            TaskStatus = "new"
	TaskStatusInProgress     TaskStatus = "in_progress"
	TaskStatusNoAnswer       TaskStatus = "no_answer"
	TaskStatusCallback       TaskStatus = "callback"
	TaskStatusContacted      TaskStatus = "contacted"
	TaskStatusBooked         TaskStatus = "booked"
	TaskStatusDeclined       TaskStatus = "declined"
	TaskStatusHandedToDoctor TaskStatus = "handed_to_doctor"
	TaskStatusCancelled      TaskStatus = "cancelled"
)

func (s TaskStatus) Active() bool {
	switch s {
	case TaskStatusNew, TaskStatusInProgress, TaskStatusNoAnswer, TaskStatusCallback:
		return true
	}
	return false
}

type DeclineReason string

const (
	DeclineReasonExpensive   DeclineReason = "expensive"
	DeclineReasonFar         DeclineReason = "far"
	DeclineReasonOtherClinic DeclineReason = "other_clinic"
	DeclineReasonNotNeeded   DeclineReason = "not_needed"
	DeclineReasonOther       DeclineReason = "other"
)

func (r DeclineReason) Valid() bool {
	switch r {
	case DeclineReasonExpensive, DeclineReasonFar, DeclineReasonOtherClinic, DeclineReasonNotNeeded, DeclineReasonOther:
		return true
	}
	return false
}

type CallOutcome string

const (
	CallOutcomeNoAnswer       CallOutcome = "no_answer"
	CallOutcomeCallback       CallOutcome = "callback"
	CallOutcomeContacted      CallOutcome = "contacted"
	CallOutcomeDeclined       CallOutcome = "declined"
	CallOutcomeBooked         CallOutcome = "booked"
	CallOutcomeHandedToDoctor CallOutcome = "handed_to_doctor"
)

type OperatorTask struct {
	ID         uuid.UUID
	CaseID     uuid.UUID
	Reason     TaskReason
	Status     TaskStatus
	Assignee   string
	DueAt      time.Time
	NextCallAt time.Time
	Attempts   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ClosedAt   time.Time
}

func NewEmergencyTask(caseID uuid.UUID, now time.Time) OperatorTask {
	return newTask(caseID, TaskReasonEmergency, UrgencyEmergency, now)
}

func NewFollowUpTask(caseID uuid.UUID, urgency Urgency, now time.Time) OperatorTask {
	return newTask(caseID, TaskReasonNoBooking, urgency, now)
}

func NewHelpRequestTask(caseID uuid.UUID, urgency Urgency, now time.Time) OperatorTask {
	return newTask(caseID, TaskReasonHelpRequest, urgency, now)
}

func newTask(caseID uuid.UUID, reason TaskReason, urgency Urgency, now time.Time) OperatorTask {
	return OperatorTask{
		CaseID:    caseID,
		Reason:    reason,
		Status:    TaskStatusNew,
		DueAt:     taskDueAt(reason, urgency, now),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t OperatorTask) NeedsDoctor() bool {
	return t.Attempts >= MaxFailedCallAttempts
}

func (t *OperatorTask) Take(actor string, now time.Time) error {
	if err := t.ensureActive(); err != nil {
		return err
	}
	t.Assignee = actor
	t.Status = TaskStatusInProgress
	t.UpdatedAt = now
	return nil
}

func (t *OperatorTask) NoAnswer(actor string, urgency Urgency, now time.Time) (bool, error) {
	if err := t.ensureActive(); err != nil {
		return false, err
	}
	t.Attempts++
	t.Assignee = actor
	t.Status = TaskStatusNoAnswer
	t.NextCallAt = now.Add(callRetryDelay(urgency))
	t.UpdatedAt = now
	return t.Attempts == MaxFailedCallAttempts, nil
}

func (t *OperatorTask) Callback(actor string, at, now time.Time) error {
	if err := t.ensureActive(); err != nil {
		return err
	}
	if !at.After(now) {
		return apperr.Invalid("call_at must be in the future")
	}
	t.Assignee = actor
	t.Status = TaskStatusCallback
	t.NextCallAt = at
	t.UpdatedAt = now
	return nil
}

func (t *OperatorTask) Close(status TaskStatus, now time.Time) error {
	if err := t.ensureActive(); err != nil {
		return err
	}
	t.Status = status
	t.ClosedAt = now
	t.UpdatedAt = now
	return nil
}

func (t OperatorTask) ensureActive() error {
	if !t.Status.Active() {
		return fmt.Errorf("%w: task is %s", apperr.ErrInvalidState, t.Status)
	}
	return nil
}

type CallAttempt struct {
	ID            uuid.UUID
	TaskID        uuid.UUID
	Outcome       CallOutcome
	DeclineReason DeclineReason
	Comment       string
	CallbackAt    time.Time
	Actor         string
	CreatedAt     time.Time
}

func NewDeclineAttempt(taskID uuid.UUID, reason DeclineReason, comment, actor string, now time.Time) (CallAttempt, error) {
	switch {
	case !reason.Valid():
		return CallAttempt{}, apperr.Invalid("reason must be one of expensive, far, other_clinic, not_needed, other")
	case reason == DeclineReasonOther && strings.TrimSpace(comment) == "":
		return CallAttempt{}, apperr.Invalid("comment is required when reason is other")
	}
	return CallAttempt{
		TaskID:        taskID,
		Outcome:       CallOutcomeDeclined,
		DeclineReason: reason,
		Comment:       comment,
		Actor:         actor,
		CreatedAt:     now,
	}, nil
}

func NewCallAttempt(taskID uuid.UUID, outcome CallOutcome, comment, actor string, now time.Time) CallAttempt {
	return CallAttempt{TaskID: taskID, Outcome: outcome, Comment: comment, Actor: actor, CreatedAt: now}
}

type TaskListItem struct {
	Task         OperatorTask
	Route        Route
	Patient      Patient
	OfferService string
}
