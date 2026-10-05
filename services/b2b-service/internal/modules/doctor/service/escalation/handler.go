package escalation

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/api"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

type Handler struct {
	tx        trm.Manager
	clock     Clock
	cases     Cases
	events    Events
	publisher Publisher
}

func NewHandler(tx trm.Manager, clock Clock, cases Cases, events Events, publisher Publisher) *Handler {
	return &Handler{tx: tx, clock: clock, cases: cases, events: events, publisher: publisher}
}

func (h *Handler) Kind() string {
	return domain.JobKindEscalateReview
}

func (h *Handler) Handle(ctx context.Context, job jobs.Job) error {
	caseID, err := uuid.Parse(job.Key)
	if err != nil {
		return fmt.Errorf("escalation job key: %w", err)
	}
	return h.tx.Do(ctx, func(ctx context.Context) error {
		c, err := h.cases.Get(ctx, caseID)
		if err != nil {
			return err
		}
		if c.Status != domain.CaseStatusInReview || c.Urgency != domain.UrgencyEmergency {
			return nil
		}
		opened, err := h.events.Exists(ctx, c.ID, domain.CaseEventOpened)
		if err != nil || opened {
			return err
		}
		return h.publisher.Publish(ctx, api.TopicReviewEscalated, c.ID.String(), api.ReviewEscalated{
			CaseID:            c.ID,
			PatientExternalID: c.Patient.ExternalID,
			EscalatedAt:       h.clock.Now(),
		})
	})
}
