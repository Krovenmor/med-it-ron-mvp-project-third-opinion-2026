package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
)

type CaseStatus string

const (
	CaseStatusInReview  CaseStatus = "in_review"
	CaseStatusConfirmed CaseStatus = "confirmed"
)

type Modality string

type Study struct {
	ID          string
	Modality    Modality
	BodySite    string
	PerformedAt time.Time
}

type Case struct {
	ID                uuid.UUID
	Patient           Patient
	Study             Study
	Conclusion        string
	Status            CaseStatus
	Urgency           Urgency
	CatalogVersion    string
	GuidelinesVersion string
	ReceivedAt        time.Time
	AssessedAt        time.Time
	ConfirmedAt       time.Time
	UpdatedAt         time.Time
}

func (c Case) ReviewDueAt() time.Time {
	return c.AssessedAt.Add(c.Urgency.ReviewSLA())
}

func (c Case) EnsureInReview() error {
	if c.Status != CaseStatusInReview {
		return fmt.Errorf("%w: case is %s, expected %s", apperr.ErrInvalidState, c.Status, CaseStatusInReview)
	}
	return nil
}

func (c *Case) ChangeUrgency(to Urgency, reason string, now time.Time) error {
	if err := c.EnsureInReview(); err != nil {
		return err
	}
	switch {
	case !to.Valid():
		return apperr.Invalid("urgency must be one of normal, planned, priority, emergency")
	case to == c.Urgency:
		return apperr.Invalid("urgency is already " + string(to))
	case to.rank() < c.Urgency.rank() && strings.TrimSpace(reason) == "":
		return apperr.Invalid("reason is required when lowering urgency")
	}
	c.Urgency = to
	c.UpdatedAt = now
	return nil
}

func (c *Case) Confirm(recs []Recommendation, now time.Time) error {
	if err := c.EnsureInReview(); err != nil {
		return err
	}
	pending := 0
	for _, r := range recs {
		if !r.Reviewed() {
			pending++
		}
	}
	if pending > 0 {
		return fmt.Errorf("%w: %d recommendations are not reviewed", apperr.ErrInvalidState, pending)
	}
	c.Status = CaseStatusConfirmed
	c.ConfirmedAt = now
	c.UpdatedAt = now
	return nil
}

type ReviewQueueItem struct {
	CaseID                  uuid.UUID
	Urgency                 Urgency
	Modality                Modality
	PerformedAt             time.Time
	ReceivedAt              time.Time
	Patient                 Patient
	RecommendationsTotal    int
	RecommendationsReviewed int
	OpenedAt                time.Time
}
