package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Urgency string

const (
	UrgencyNormal    Urgency = "normal"
	UrgencyPlanned   Urgency = "planned"
	UrgencyPriority  Urgency = "priority"
	UrgencyEmergency Urgency = "emergency"
)

func (u Urgency) Valid() bool {
	switch u {
	case UrgencyNormal, UrgencyPlanned, UrgencyPriority, UrgencyEmergency:
		return true
	}
	return false
}

type Importance string

const (
	ImportanceHigh Importance = "high"
	ImportanceLow  Importance = "low"
)

func (i Importance) Valid() bool {
	return i == ImportanceHigh || i == ImportanceLow
}

type Service struct {
	Code string
	Name string
}

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

type Recommendation struct {
	ServiceCode   string
	ServiceName   string
	Importance    Importance
	Rationale     string
	GuidelineRef  string
	PatientText   string
	AlreadyBooked bool
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
