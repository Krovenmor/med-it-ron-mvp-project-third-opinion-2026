package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AssessmentRequest struct {
	CaseID                 uuid.UUID
	Modality               Modality
	BodySite               string
	PerformedAt            time.Time
	Conclusion             string
	PatientAge             int
	PatientSex             Sex
	CurrentRecommendations []Service
	History                PatientHistory
}

type Assessment struct {
	Urgency           Urgency
	CatalogVersion    string
	GuidelinesVersion string
	Recommendations   []Recommendation
}

func (a Assessment) Validate() error {
	switch {
	case !a.Urgency.Valid():
		return fmt.Errorf("unknown urgency %q", a.Urgency)
	case a.CatalogVersion == "":
		return errors.New("catalog_version is empty")
	case a.GuidelinesVersion == "":
		return errors.New("guidelines_version is empty")
	}
	for i, r := range a.Recommendations {
		switch {
		case r.ServiceName == "":
			return fmt.Errorf("recommendation %d: service_name is empty", i)
		case r.PatientText == "":
			return fmt.Errorf("recommendation %d: patient_text is empty", i)
		case !r.Importance.Valid():
			return fmt.Errorf("recommendation %d: unknown importance %q", i, r.Importance)
		}
	}
	return nil
}

func (a Assessment) RecommendationsFor(caseID uuid.UUID, now time.Time) []Recommendation {
	recs := make([]Recommendation, len(a.Recommendations))
	for i, r := range a.Recommendations {
		r.CaseID = caseID
		r.Position = i + 1
		r.Source = RecommendationSourceAI
		r.CreatedAt = now
		recs[i] = r
	}
	return recs
}
