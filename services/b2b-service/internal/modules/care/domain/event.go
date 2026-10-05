package domain

import (
	"time"

	"github.com/google/uuid"
)

type CaseEventType string

const (
	CaseEventStatusChanged          CaseEventType = "status_changed"
	CaseEventRecommendationBooked   CaseEventType = "recommendation_booked"
	CaseEventTaskCreated            CaseEventType = "operator_task_created"
	CaseEventCallAttempt            CaseEventType = "call_attempt"
	CaseEventTaskClosed             CaseEventType = "operator_task_closed"
	CaseEventRecommendationDeclined CaseEventType = "recommendation_declined"
)

const ActorSystem = "system"

type CaseEvent struct {
	CaseID     uuid.UUID
	Type       CaseEventType
	FromStatus RouteStatus
	ToStatus   RouteStatus
	Actor      string
	Payload    map[string]string
	OccurredAt time.Time
}

func StatusChanged(caseID uuid.UUID, from, to RouteStatus, actor string, at time.Time) CaseEvent {
	return CaseEvent{CaseID: caseID, Type: CaseEventStatusChanged, FromStatus: from, ToStatus: to, Actor: actor, OccurredAt: at}
}

func RecommendationBooked(b Booking) CaseEvent {
	return CaseEvent{
		CaseID: b.CaseID,
		Type:   CaseEventRecommendationBooked,
		Actor:  b.Channel.Actor(),
		Payload: map[string]string{
			"recommendation_id": b.RecommendationID.String(),
			"appointment_id":    b.AppointmentID,
			"scheduled_at":      b.ScheduledAt.Format(time.RFC3339),
			"channel":           string(b.Channel),
		},
		OccurredAt: b.CreatedAt,
	}
}

func RecommendationDeclined(i PlanItem, actor string) CaseEvent {
	return CaseEvent{
		CaseID: i.CaseID,
		Type:   CaseEventRecommendationDeclined,
		Actor:  actor,
		Payload: map[string]string{
			"recommendation_id": i.ID.String(),
			"reason":            string(i.DeclineReason),
			"comment":           i.DeclineComment,
		},
		OccurredAt: i.DeclinedAt,
	}
}

func OperatorTaskCreated(t OperatorTask) CaseEvent {
	return CaseEvent{
		CaseID: t.CaseID,
		Type:   CaseEventTaskCreated,
		Actor:  ActorSystem,
		Payload: map[string]string{
			"task_id": t.ID.String(),
			"reason":  string(t.Reason),
			"due_at":  t.DueAt.Format(time.RFC3339),
		},
		OccurredAt: t.CreatedAt,
	}
}

func CallAttempted(caseID uuid.UUID, a CallAttempt) CaseEvent {
	return CaseEvent{
		CaseID: caseID,
		Type:   CaseEventCallAttempt,
		Actor:  a.Actor,
		Payload: map[string]string{
			"task_id":        a.TaskID.String(),
			"outcome":        string(a.Outcome),
			"decline_reason": string(a.DeclineReason),
			"comment":        a.Comment,
		},
		OccurredAt: a.CreatedAt,
	}
}

func OperatorTaskClosed(t OperatorTask, actor string) CaseEvent {
	return CaseEvent{
		CaseID: t.CaseID,
		Type:   CaseEventTaskClosed,
		Actor:  actor,
		Payload: map[string]string{
			"task_id": t.ID.String(),
			"status":  string(t.Status),
		},
		OccurredAt: t.ClosedAt,
	}
}
