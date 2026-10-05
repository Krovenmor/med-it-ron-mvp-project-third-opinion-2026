package domain

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

type RouteStatus string

const (
	RouteStatusReceived    RouteStatus = "received"
	RouteStatusInReview    RouteStatus = "in_review"
	RouteStatusConfirmed   RouteStatus = "confirmed"
	RouteStatusNotified    RouteStatus = "notified"
	RouteStatusBooked      RouteStatus = "booked"
	RouteStatusCompleted   RouteStatus = "completed"
	RouteStatusDeclined    RouteStatus = "declined"
	RouteStatusUnreachable RouteStatus = "unreachable"
)

type Route struct {
	CaseID      uuid.UUID
	PatientID   uuid.UUID
	Status      RouteStatus
	Urgency     Urgency
	Modality    Modality
	PerformedAt time.Time
	ReceivedAt  time.Time
	AssessedAt  time.Time
	ReviewDueAt time.Time
	ConfirmedAt time.Time
	UpdatedAt   time.Time
}

type ReviewStart struct {
	CaseID      uuid.UUID
	PatientID   uuid.UUID
	Urgency     Urgency
	Modality    Modality
	PerformedAt time.Time
	ReceivedAt  time.Time
	AssessedAt  time.Time
	ReviewDueAt time.Time
}

func NewRoute(caseID, patientID uuid.UUID, modality Modality, performedAt, receivedAt time.Time) Route {
	return Route{
		CaseID:      caseID,
		PatientID:   patientID,
		Status:      RouteStatusReceived,
		Modality:    modality,
		PerformedAt: performedAt,
		ReceivedAt:  receivedAt,
		UpdatedAt:   receivedAt,
	}
}

func (r *Route) StartReview(s ReviewStart) bool {
	r.CaseID, r.PatientID = s.CaseID, s.PatientID
	r.Modality, r.PerformedAt, r.ReceivedAt = s.Modality, s.PerformedAt, s.ReceivedAt
	r.Urgency, r.AssessedAt, r.ReviewDueAt = s.Urgency, s.AssessedAt, s.ReviewDueAt
	r.UpdatedAt = s.AssessedAt
	if r.Status != "" && r.Status != RouteStatusReceived {
		return false
	}
	r.Status = RouteStatusInReview
	return true
}

func (r *Route) ChangeUrgency(to Urgency, reviewDueAt, now time.Time) {
	r.Urgency = to
	r.ReviewDueAt = reviewDueAt
	r.UpdatedAt = now
}

func (r *Route) Confirm(urgency Urgency, at time.Time) bool {
	r.Urgency = urgency
	if !r.moveFrom([]RouteStatus{RouteStatusInReview}, RouteStatusConfirmed, at) {
		return false
	}
	r.ConfirmedAt = at
	return true
}

func (r Route) AwaitingBooking() bool {
	return r.Status == RouteStatusConfirmed || r.Status == RouteStatusNotified
}

func (r Route) Confirmed() bool {
	switch r.Status {
	case RouteStatusConfirmed, RouteStatusNotified, RouteStatusBooked, RouteStatusCompleted, RouteStatusDeclined, RouteStatusUnreachable:
		return true
	}
	return false
}

func (r Route) NeedsUrgentContact() bool {
	if r.Urgency != UrgencyEmergency {
		return false
	}
	switch r.Status {
	case RouteStatusInReview, RouteStatusConfirmed, RouteStatusNotified:
		return true
	}
	return false
}

func (r Route) NoFindings() bool {
	return r.Urgency == UrgencyNormal && r.Confirmed()
}

func (r Route) InReview() bool {
	return r.Status == RouteStatusInReview
}

func (r Route) BookBy() time.Time {
	start := r.ConfirmedAt
	if start.IsZero() {
		start = r.ReceivedAt
	}
	return start.Add(bookByOffset(r.Urgency))
}

func (r *Route) MarkNotified(now time.Time) bool {
	return r.moveFrom([]RouteStatus{RouteStatusConfirmed}, RouteStatusNotified, now)
}

func (r *Route) Decline(now time.Time) bool {
	return r.moveFrom([]RouteStatus{RouteStatusConfirmed, RouteStatusNotified, RouteStatusUnreachable}, RouteStatusDeclined, now)
}

func (r *Route) MarkUnreachable(now time.Time) bool {
	return r.moveFrom([]RouteStatus{RouteStatusConfirmed, RouteStatusNotified}, RouteStatusUnreachable, now)
}

func (r *Route) Complete(now time.Time) bool {
	return r.moveFrom([]RouteStatus{RouteStatusBooked}, RouteStatusCompleted, now)
}

func (r Route) EnsureBookable() error {
	switch r.Status {
	case RouteStatusConfirmed, RouteStatusNotified, RouteStatusBooked, RouteStatusUnreachable, RouteStatusDeclined:
		return nil
	}
	return fmt.Errorf("%w: case is %s and cannot be booked", apperr.ErrInvalidState, r.Status)
}

func (r *Route) Book(now time.Time) (bool, error) {
	switch r.Status {
	case RouteStatusBooked:
		return false, nil
	case RouteStatusConfirmed, RouteStatusNotified, RouteStatusUnreachable, RouteStatusDeclined:
		r.Status = RouteStatusBooked
		r.UpdatedAt = now
		return true, nil
	}
	return false, fmt.Errorf("%w: case is %s and cannot be booked", apperr.ErrInvalidState, r.Status)
}

func (r *Route) moveFrom(allowed []RouteStatus, to RouteStatus, now time.Time) bool {
	if !slices.Contains(allowed, r.Status) {
		return false
	}
	r.Status = to
	r.UpdatedAt = now
	return true
}
