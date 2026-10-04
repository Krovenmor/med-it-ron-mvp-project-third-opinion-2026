package http

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/cases"
)

type reportRequest struct {
	SourceSystem string         `json:"source_system"`
	Study        studyRequest   `json:"study"`
	Conclusion   string         `json:"conclusion"`
	Patient      patientRequest `json:"patient"`
}

type studyRequest struct {
	ID          string `json:"id"`
	Modality    string `json:"modality"`
	BodySite    string `json:"body_site"`
	PerformedAt string `json:"performed_at"`
}

type patientRequest struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func (r reportRequest) toDomain() (domain.Report, error) {
	study, err := r.Study.toDomain()
	if err != nil {
		return domain.Report{}, err
	}
	patient, err := r.Patient.toDomain(r.SourceSystem)
	if err != nil {
		return domain.Report{}, err
	}
	return domain.Report{
		SourceSystem: r.SourceSystem,
		Study:        study,
		Conclusion:   r.Conclusion,
		Patient:      patient,
	}, nil
}

func (s studyRequest) toDomain() (domain.Study, error) {
	performedAt, err := parseOptionalTime("study.performed_at", s.PerformedAt, time.RFC3339, "an RFC 3339 date-time")
	if err != nil {
		return domain.Study{}, err
	}
	return domain.Study{
		ID:          s.ID,
		Modality:    domain.Modality(s.Modality),
		BodySite:    s.BodySite,
		PerformedAt: performedAt,
	}, nil
}

func (p patientRequest) toDomain(sourceSystem string) (domain.Patient, error) {
	birthDate, err := parseOptionalTime("patient.birth_date", p.BirthDate, time.DateOnly, "YYYY-MM-DD")
	if err != nil {
		return domain.Patient{}, err
	}
	return domain.Patient{
		SourceSystem: sourceSystem,
		ExternalID:   p.ID,
		FullName:     p.FullName,
		BirthDate:    birthDate,
		Sex:          domain.Sex(p.Sex),
		Phone:        p.Phone,
		Email:        p.Email,
	}, nil
}

func parseOptionalTime(field, value, layout, format string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s must be %s", domain.ErrInvalidInput, field, format)
	}
	return t, nil
}

type caseRefResponse struct {
	CaseID uuid.UUID `json:"case_id"`
	Status string    `json:"status"`
}

type caseResponse struct {
	ID                uuid.UUID                `json:"id"`
	PatientID         uuid.UUID                `json:"patient_id"`
	SourceSystem      string                   `json:"source_system"`
	StudyID           string                   `json:"study_id"`
	Modality          string                   `json:"modality"`
	BodySite          string                   `json:"body_site"`
	PerformedAt       time.Time                `json:"performed_at"`
	Conclusion        string                   `json:"conclusion"`
	Status            string                   `json:"status"`
	Urgency           string                   `json:"urgency,omitempty"`
	CatalogVersion    string                   `json:"catalog_version,omitempty"`
	GuidelinesVersion string                   `json:"guidelines_version,omitempty"`
	ReceivedAt        time.Time                `json:"received_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	Recommendations   []recommendationResponse `json:"recommendations"`
}

type recommendationResponse struct {
	ID            uuid.UUID `json:"id"`
	Position      int       `json:"position"`
	Source        string    `json:"source"`
	ServiceCode   string    `json:"service_code"`
	ServiceName   string    `json:"service_name"`
	Importance    string    `json:"importance"`
	Rationale     string    `json:"rationale"`
	GuidelineRef  string    `json:"guideline_ref"`
	PatientText   string    `json:"patient_text"`
	AlreadyBooked bool      `json:"already_booked"`
}

func newCaseResponse(d cases.Details) caseResponse {
	c := d.Case
	recs := make([]recommendationResponse, 0, len(d.Recommendations))
	for _, r := range d.Recommendations {
		recs = append(recs, recommendationResponse{
			ID:            r.ID,
			Position:      r.Position,
			Source:        string(r.Source),
			ServiceCode:   r.ServiceCode,
			ServiceName:   r.ServiceName,
			Importance:    string(r.Importance),
			Rationale:     r.Rationale,
			GuidelineRef:  r.GuidelineRef,
			PatientText:   r.PatientText,
			AlreadyBooked: r.AlreadyBooked,
		})
	}
	return caseResponse{
		ID:                c.ID,
		PatientID:         c.PatientID,
		SourceSystem:      c.SourceSystem,
		StudyID:           c.Study.ID,
		Modality:          string(c.Study.Modality),
		BodySite:          c.Study.BodySite,
		PerformedAt:       c.Study.PerformedAt,
		Conclusion:        c.Conclusion,
		Status:            string(c.Status),
		Urgency:           string(c.Urgency),
		CatalogVersion:    c.CatalogVersion,
		GuidelinesVersion: c.GuidelinesVersion,
		ReceivedAt:        c.ReceivedAt,
		UpdatedAt:         c.UpdatedAt,
		Recommendations:   recs,
	}
}

type advanceClockRequest struct {
	By string `json:"by"`
}

type clockResponse struct {
	Now time.Time `json:"now"`
}

type errorResponse struct {
	Error string `json:"error"`
}
