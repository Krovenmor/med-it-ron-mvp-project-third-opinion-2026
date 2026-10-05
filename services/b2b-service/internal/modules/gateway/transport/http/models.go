package http

import (
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
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

type reportRefResponse struct {
	CaseID uuid.UUID `json:"case_id"`
	Status string    `json:"status"`
}

type reportStatusResponse struct {
	CaseID     uuid.UUID  `json:"case_id"`
	StudyID    string     `json:"study_id"`
	Status     string     `json:"status"`
	ReceivedAt time.Time  `json:"received_at"`
	AssessedAt *time.Time `json:"assessed_at"`
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
	performedAt, err := httpx.ParseOptionalTime("study.performed_at", s.PerformedAt, time.RFC3339, "an RFC 3339 date-time")
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
	birthDate, err := httpx.ParseOptionalTime("patient.birth_date", p.BirthDate, time.DateOnly, "YYYY-MM-DD")
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

func newReportRefResponse(i domain.Intake) reportRefResponse {
	return reportRefResponse{CaseID: i.CaseID, Status: string(i.Status)}
}

func newReportStatusResponse(i domain.Intake) reportStatusResponse {
	return reportStatusResponse{
		CaseID:     i.CaseID,
		StudyID:    i.Study.ID,
		Status:     string(i.Status),
		ReceivedAt: i.ReceivedAt,
		AssessedAt: httpx.OptionalTime(i.AssessedAt),
	}
}
