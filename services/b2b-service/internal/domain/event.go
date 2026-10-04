package domain

import (
	"time"

	"github.com/google/uuid"
)

type CaseEventType string

const CaseEventStatusChanged CaseEventType = "status_changed"

const ActorSystem = "system"

type CaseEvent struct {
	CaseID     uuid.UUID
	Type       CaseEventType
	FromStatus CaseStatus
	ToStatus   CaseStatus
	Actor      string
	OccurredAt time.Time
}

func StatusChanged(caseID uuid.UUID, from, to CaseStatus, actor string, at time.Time) CaseEvent {
	return CaseEvent{
		CaseID:     caseID,
		Type:       CaseEventStatusChanged,
		FromStatus: from,
		ToStatus:   to,
		Actor:      actor,
		OccurredAt: at,
	}
}
