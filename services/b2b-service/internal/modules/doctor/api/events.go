package api

import (
	"time"

	"github.com/google/uuid"
)

const (
	TopicReviewStarted   = "doctor.review_started"
	TopicUrgencyChanged  = "doctor.urgency_changed"
	TopicCaseConfirmed   = "doctor.case_confirmed"
	TopicReviewEscalated = "doctor.review_escalated"
)

type ReviewStarted struct {
	CaseID      uuid.UUID `json:"case_id"`
	PatientID   uuid.UUID `json:"patient_id"`
	Urgency     string    `json:"urgency"`
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
	ReceivedAt  time.Time `json:"received_at"`
	AssessedAt  time.Time `json:"assessed_at"`
	ReviewDueAt time.Time `json:"review_due_at"`
}

type UrgencyChanged struct {
	CaseID      uuid.UUID `json:"case_id"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Reason      string    `json:"reason"`
	Actor       string    `json:"actor"`
	ReviewDueAt time.Time `json:"review_due_at"`
	ChangedAt   time.Time `json:"changed_at"`
}

type CaseConfirmed struct {
	CaseID          uuid.UUID                 `json:"case_id"`
	Urgency         string                    `json:"urgency"`
	Actor           string                    `json:"actor"`
	ConfirmedAt     time.Time                 `json:"confirmed_at"`
	Recommendations []ConfirmedRecommendation `json:"recommendations"`
}

type ConfirmedRecommendation struct {
	ID          uuid.UUID `json:"id"`
	Position    int       `json:"position"`
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	PatientText string    `json:"patient_text"`
	Mark        string    `json:"mark"`
}

type ReviewEscalated struct {
	CaseID            uuid.UUID `json:"case_id"`
	PatientExternalID string    `json:"patient_external_id"`
	EscalatedAt       time.Time `json:"escalated_at"`
}
