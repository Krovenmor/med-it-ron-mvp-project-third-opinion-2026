//go:build e2e

package e2e

import (
	"time"

	"github.com/google/uuid"
)

type report struct {
	SourceSystem string        `json:"source_system"`
	Study        reportStudy   `json:"study"`
	Conclusion   string        `json:"conclusion"`
	Patient      reportPatient `json:"patient"`
}

type reportStudy struct {
	ID          string `json:"id"`
	Modality    string `json:"modality"`
	BodySite    string `json:"body_site"`
	PerformedAt string `json:"performed_at"`
}

type reportPatient struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func newReport() report {
	id := uuid.NewString()
	return report{
		SourceSystem: "e2e",
		Study: reportStudy{
			ID:          "STUDY-" + id,
			Modality:    "CT",
			BodySite:    "chest",
			PerformedAt: "2026-10-01T09:00:00Z",
		},
		Conclusion: "Солидный очаг 9 мм, исследование " + id,
		Patient: reportPatient{
			ID:        "PATIENT-" + id,
			FullName:  "Пациент " + id,
			BirthDate: "1974-03-12",
			Sex:       "female",
			Phone:     "+7-" + id,
			Email:     id + "@example.com",
		},
	}
}

type caseRef struct {
	CaseID string `json:"case_id"`
	Status string `json:"status"`
}

type caseView struct {
	ID                string               `json:"id"`
	StudyID           string               `json:"study_id"`
	Status            string               `json:"status"`
	Urgency           string               `json:"urgency"`
	CatalogVersion    string               `json:"catalog_version"`
	GuidelinesVersion string               `json:"guidelines_version"`
	ReceivedAt        time.Time            `json:"received_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	Patient           patientView          `json:"patient"`
	Recommendations   []caseRecommendation `json:"recommendations"`
}

type patientView struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
}

type caseRecommendation struct {
	recommendation
	ID       string      `json:"id"`
	Position int         `json:"position"`
	Source   string      `json:"source"`
	Review   *reviewView `json:"review"`
}

type reviewView struct {
	Mark          string    `json:"mark"`
	RejectReason  string    `json:"reject_reason"`
	RejectComment string    `json:"reject_comment"`
	ReviewedBy    string    `json:"reviewed_by"`
	ReviewedAt    time.Time `json:"reviewed_at"`
}

type reviewQueueView struct {
	Cases []queueCaseView `json:"cases"`
}

type queueCaseView struct {
	CaseID                  string      `json:"case_id"`
	Urgency                 string      `json:"urgency"`
	Modality                string      `json:"modality"`
	Patient                 patientView `json:"patient"`
	RecommendationsTotal    int         `json:"recommendations_total"`
	RecommendationsReviewed int         `json:"recommendations_reviewed"`
	OpenedAt                time.Time   `json:"opened_at"`
}

type reviewRequest struct {
	Mark          string  `json:"mark"`
	RejectReason  string  `json:"reject_reason,omitempty"`
	RejectComment string  `json:"reject_comment,omitempty"`
	PatientText   *string `json:"patient_text,omitempty"`
}

type addRecommendationRequest struct {
	ServiceCode string `json:"service_code,omitempty"`
	ServiceName string `json:"service_name"`
	Rationale   string `json:"rationale,omitempty"`
	PatientText string `json:"patient_text"`
	Mark        string `json:"mark"`
}

type urgencyRequest struct {
	Urgency string `json:"urgency"`
	Reason  string `json:"reason,omitempty"`
}

type planView struct {
	Patient planPatientView `json:"patient"`
	Cases   []planCaseView  `json:"cases"`
}

type planPatientView struct {
	FullName string `json:"full_name"`
}

type planCaseView struct {
	CaseID          string               `json:"case_id"`
	Status          string               `json:"status"`
	UrgentContact   bool                 `json:"urgent_contact"`
	Study           planStudyView        `json:"study"`
	Recommendations []planRecommendation `json:"recommendations"`
}

type planStudyView struct {
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
}

type planRecommendation struct {
	ID          string `json:"id"`
	ServiceCode string `json:"service_code"`
	ServiceName string `json:"service_name"`
	PatientText string `json:"patient_text"`
	Mark        string `json:"mark"`
}

type bookingRequest struct {
	RecommendationID string `json:"recommendation_id"`
	AppointmentID    string `json:"appointment_id"`
	ScheduledAt      string `json:"scheduled_at"`
	Channel          string `json:"channel"`
}

type bookingView struct {
	BookingID        string `json:"booking_id"`
	CaseID           string `json:"case_id"`
	RecommendationID string `json:"recommendation_id"`
	AppointmentID    string `json:"appointment_id"`
	CaseStatus       string `json:"case_status"`
}

type caseEventView struct {
	Type    string            `db:"type"`
	Actor   string            `db:"actor"`
	Payload map[string]string `db:"payload"`
}

type assessment struct {
	Urgency           string           `json:"urgency"`
	CatalogVersion    string           `json:"catalog_version"`
	GuidelinesVersion string           `json:"guidelines_version"`
	Recommendations   []recommendation `json:"recommendations"`
}

type recommendation struct {
	ServiceCode   string `json:"service_code"`
	ServiceName   string `json:"service_name"`
	Importance    string `json:"importance"`
	Rationale     string `json:"rationale"`
	GuidelineRef  string `json:"guideline_ref"`
	PatientText   string `json:"patient_text"`
	AlreadyBooked bool   `json:"already_booked"`
}

func validAssessment() assessment {
	return assessment{
		Urgency:           "priority",
		CatalogVersion:    "catalog-e2e",
		GuidelinesVersion: "guidelines-e2e",
		Recommendations: []recommendation{
			{
				ServiceCode:  "PULM-CONSULT",
				ServiceName:  "Консультация пульмонолога",
				Importance:   "high",
				Rationale:    "Солидный очаг более 8 мм",
				GuidelineRef: "КР МЗ РФ",
				PatientText:  "Рекомендуем консультацию пульмонолога",
			},
			{
				ServiceName:   "Консультация торакального хирурга",
				Importance:    "low",
				PatientText:   "Рекомендуем консультацию хирурга в другой клинике",
				AlreadyBooked: true,
			},
		},
	}
}

type history struct {
	Visits       []visit       `json:"visits"`
	Appointments []appointment `json:"appointments"`
}

type visit struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	VisitedAt   time.Time `json:"visited_at"`
}

type appointment struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type aiRequest struct {
	CaseID       string        `json:"case_id"`
	Study        aiStudy       `json:"study"`
	Conclusion   string        `json:"conclusion"`
	Patient      aiPatient     `json:"patient"`
	Visits       []visit       `json:"visits"`
	Appointments []appointment `json:"appointments"`
}

type aiStudy struct {
	Modality    string    `json:"modality"`
	BodySite    string    `json:"body_site"`
	PerformedAt time.Time `json:"performed_at"`
}

type aiPatient struct {
	Age int    `json:"age"`
	Sex string `json:"sex"`
}

type clockView struct {
	Now time.Time `json:"now"`
}

type jobState struct {
	State    string `db:"state"`
	Attempts int    `db:"attempts"`
}
