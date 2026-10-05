package domain

import (
	"fmt"
	"time"
)

type Plan struct {
	Now         time.Time
	PatientName string
	Cases       []PlanCase
}

type PlanCase struct {
	ID              string
	Status          string
	UrgentContact   bool
	NoFindings      bool
	Modality        string
	PerformedAt     time.Time
	ConfirmedAt     time.Time
	BookBy          time.Time
	Recommendations []Recommendation
}

func (c PlanCase) InReview() bool {
	return c.Status == "in_review"
}

func (c PlanCase) AssignedAt() time.Time {
	if c.ConfirmedAt.IsZero() {
		return c.PerformedAt
	}
	return c.ConfirmedAt
}

type Recommendation struct {
	ID            string
	CaseID        string
	ServiceCode   string
	ServiceName   string
	PatientText   string
	Mark          string
	Declined      bool
	DeclineReason string
}

func (r Recommendation) InClinic() bool {
	return r.ServiceCode != ""
}

func (r Recommendation) EnsureBookable() error {
	if !r.InClinic() {
		return fmt.Errorf("%w: service is not available in the clinic", ErrInvalidState)
	}
	return nil
}

func (p Plan) Recommendation(id string) (Recommendation, error) {
	for _, c := range p.Cases {
		for _, r := range c.Recommendations {
			if r.ID == id {
				return r, nil
			}
		}
	}
	return Recommendation{}, fmt.Errorf("recommendation %s: %w", id, ErrNotFound)
}

type Decline struct {
	RecommendationID string
	Reason           string
	Comment          string
}
