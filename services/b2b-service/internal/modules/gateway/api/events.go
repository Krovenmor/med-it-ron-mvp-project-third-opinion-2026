package api

import (
	"time"

	"github.com/google/uuid"
)

const (
	TopicPatientRegistered = "gateway.patient_registered"
	TopicReportReceived    = "gateway.report_received"
	TopicCaseAssessed      = "gateway.case_assessed"
)

type PatientRegistered struct {
	PatientID    uuid.UUID `json:"patient_id"`
	SourceSystem string    `json:"source_system"`
	ExternalID   string    `json:"external_id"`
	FullName     string    `json:"full_name"`
	BirthDate    time.Time `json:"birth_date"`
	Sex          string    `json:"sex"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
}

type ReportReceived struct {
	CaseID      uuid.UUID `json:"case_id"`
	PatientID   uuid.UUID `json:"patient_id"`
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
	ReceivedAt  time.Time `json:"received_at"`
}

type CaseAssessed struct {
	CaseID            uuid.UUID                `json:"case_id"`
	Patient           AssessedPatient          `json:"patient"`
	Study             AssessedStudy            `json:"study"`
	Conclusion        string                   `json:"conclusion"`
	ReceivedAt        time.Time                `json:"received_at"`
	AssessedAt        time.Time                `json:"assessed_at"`
	Urgency           string                   `json:"urgency"`
	CatalogVersion    string                   `json:"catalog_version"`
	GuidelinesVersion string                   `json:"guidelines_version"`
	Recommendations   []AssessedRecommendation `json:"recommendations"`
}

type AssessedPatient struct {
	ID           uuid.UUID `json:"id"`
	SourceSystem string    `json:"source_system"`
	ExternalID   string    `json:"external_id"`
	BirthDate    time.Time `json:"birth_date"`
	Sex          string    `json:"sex"`
}

type AssessedStudy struct {
	ID          string    `json:"id"`
	Modality    string    `json:"modality"`
	BodySite    string    `json:"body_site"`
	PerformedAt time.Time `json:"performed_at"`
}

type AssessedRecommendation struct {
	ServiceCode   string `json:"service_code"`
	ServiceName   string `json:"service_name"`
	Importance    string `json:"importance"`
	Rationale     string `json:"rationale"`
	GuidelineRef  string `json:"guideline_ref"`
	PatientText   string `json:"patient_text"`
	AlreadyBooked bool   `json:"already_booked"`
}
