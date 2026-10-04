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
	Recommendations   []caseRecommendation `json:"recommendations"`
}

type caseRecommendation struct {
	recommendation
	Position int    `json:"position"`
	Source   string `json:"source"`
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
