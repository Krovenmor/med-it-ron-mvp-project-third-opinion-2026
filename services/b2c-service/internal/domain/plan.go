package domain

import (
	"fmt"
	"time"
)

type Plan struct {
	PatientName string
	Cases       []PlanCase
}

type PlanCase struct {
	ID              string
	Status          string
	UrgentContact   bool
	Modality        string
	PerformedAt     time.Time
	Recommendations []Recommendation
}

type Recommendation struct {
	ID          string
	CaseID      string
	ServiceCode string
	ServiceName string
	PatientText string
	Mark        string
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
