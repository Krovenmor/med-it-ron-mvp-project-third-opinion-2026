package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

type Mark string

const (
	MarkCritical Mark = "critical"
	MarkMinor    Mark = "minor"
)

type PlanItem struct {
	ID             uuid.UUID
	CaseID         uuid.UUID
	Position       int
	ServiceCode    string
	ServiceName    string
	PatientText    string
	Mark           Mark
	DeclinedAt     time.Time
	DeclineReason  DeclineReason
	DeclineComment string
}

func (i PlanItem) Declined() bool {
	return !i.DeclinedAt.IsZero()
}

func (i *PlanItem) Decline(reason DeclineReason, comment string, now time.Time) error {
	switch {
	case i.Declined():
		return fmt.Errorf("%w: recommendation is already declined", apperr.ErrInvalidState)
	case !reason.Valid():
		return apperr.Invalid("reason must be one of expensive, far, other_clinic, not_needed, other")
	case reason == DeclineReasonOther && strings.TrimSpace(comment) == "":
		return apperr.Invalid("comment is required when reason is other")
	}
	i.DeclinedAt = now
	i.DeclineReason = reason
	i.DeclineComment = comment
	return nil
}

func (i PlanItem) EnsureBookable() error {
	if i.ServiceCode == "" {
		return fmt.Errorf("%w: service is not available in the clinic", apperr.ErrInvalidState)
	}
	return nil
}

type Service struct {
	Code string
	Name string
}

type PatientPlan struct {
	Now     time.Time
	Patient Patient
	Cases   []PlanCase
}

type PlanCase struct {
	Route         Route
	UrgentContact bool
	Items         []PlanItem
}
