package domain

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type CaseEventType string

const (
	CaseEventStatusChanged            CaseEventType = "status_changed"
	CaseEventOpened                   CaseEventType = "case_opened"
	CaseEventRecommendationReviewed   CaseEventType = "recommendation_reviewed"
	CaseEventRecommendationTextEdited CaseEventType = "recommendation_text_edited"
	CaseEventRecommendationAdded      CaseEventType = "recommendation_added"
	CaseEventUrgencyChanged           CaseEventType = "urgency_changed"
)

const ActorSystem = "system"

type CaseEvent struct {
	CaseID     uuid.UUID
	Type       CaseEventType
	FromStatus CaseStatus
	ToStatus   CaseStatus
	Actor      string
	Payload    map[string]string
	OccurredAt time.Time
}

func StatusChanged(caseID uuid.UUID, from, to CaseStatus, actor string, at time.Time) CaseEvent {
	return CaseEvent{CaseID: caseID, Type: CaseEventStatusChanged, FromStatus: from, ToStatus: to, Actor: actor, OccurredAt: at}
}

func CaseOpened(caseID uuid.UUID, actor string, at time.Time) CaseEvent {
	return CaseEvent{CaseID: caseID, Type: CaseEventOpened, Actor: actor, OccurredAt: at}
}

func RecommendationReviewed(rec Recommendation, patientTextEdited bool) CaseEvent {
	return CaseEvent{
		CaseID: rec.CaseID,
		Type:   CaseEventRecommendationReviewed,
		Actor:  rec.Review.ReviewedBy,
		Payload: map[string]string{
			"recommendation_id":   rec.ID.String(),
			"mark":                string(rec.Review.Mark),
			"reject_reason":       string(rec.Review.RejectReason),
			"reject_comment":      rec.Review.RejectComment,
			"patient_text_edited": strconv.FormatBool(patientTextEdited),
		},
		OccurredAt: rec.Review.ReviewedAt,
	}
}

func RecommendationTextEdited(rec Recommendation, actor string, at time.Time) CaseEvent {
	return CaseEvent{
		CaseID:     rec.CaseID,
		Type:       CaseEventRecommendationTextEdited,
		Actor:      actor,
		Payload:    map[string]string{"recommendation_id": rec.ID.String()},
		OccurredAt: at,
	}
}

func RecommendationAdded(rec Recommendation) CaseEvent {
	return CaseEvent{
		CaseID: rec.CaseID,
		Type:   CaseEventRecommendationAdded,
		Actor:  rec.Review.ReviewedBy,
		Payload: map[string]string{
			"recommendation_id": rec.ID.String(),
			"mark":              string(rec.Review.Mark),
		},
		OccurredAt: rec.CreatedAt,
	}
}

func UrgencyChanged(caseID uuid.UUID, from, to Urgency, reason, actor string, at time.Time) CaseEvent {
	return CaseEvent{
		CaseID: caseID,
		Type:   CaseEventUrgencyChanged,
		Actor:  actor,
		Payload: map[string]string{
			"from":   string(from),
			"to":     string(to),
			"reason": reason,
		},
		OccurredAt: at,
	}
}
