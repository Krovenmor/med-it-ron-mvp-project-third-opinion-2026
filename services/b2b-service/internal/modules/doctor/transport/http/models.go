package http

import (
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/cases"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/service/review"
	gatewayapi "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/api"
)

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
	Urgency           string                   `json:"urgency"`
	CatalogVersion    string                   `json:"catalog_version"`
	GuidelinesVersion string                   `json:"guidelines_version"`
	ReceivedAt        time.Time                `json:"received_at"`
	AssessedAt        time.Time                `json:"assessed_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	Patient           patientResponse          `json:"patient"`
	Recommendations   []recommendationResponse `json:"recommendations"`
}

type patientResponse struct {
	ID        string `json:"id"`
	BirthDate string `json:"birth_date"`
	Sex       string `json:"sex"`
}

type recommendationResponse struct {
	ID            uuid.UUID       `json:"id"`
	Position      int             `json:"position"`
	Source        string          `json:"source"`
	ServiceCode   string          `json:"service_code"`
	ServiceName   string          `json:"service_name"`
	Importance    string          `json:"importance"`
	Rationale     string          `json:"rationale"`
	GuidelineRef  string          `json:"guideline_ref"`
	PatientText   string          `json:"patient_text"`
	AlreadyBooked bool            `json:"already_booked"`
	Review        *reviewResponse `json:"review"`
}

type reviewResponse struct {
	Mark          string    `json:"mark"`
	RejectReason  string    `json:"reject_reason,omitempty"`
	RejectComment string    `json:"reject_comment,omitempty"`
	ReviewedBy    string    `json:"reviewed_by"`
	ReviewedAt    time.Time `json:"reviewed_at"`
}

type reviewQueueResponse struct {
	Cases []queueCaseResponse `json:"cases"`
}

type queueCaseResponse struct {
	CaseID                  uuid.UUID       `json:"case_id"`
	Urgency                 string          `json:"urgency"`
	Modality                string          `json:"modality"`
	PerformedAt             time.Time       `json:"performed_at"`
	ReceivedAt              time.Time       `json:"received_at"`
	Patient                 patientResponse `json:"patient"`
	RecommendationsTotal    int             `json:"recommendations_total"`
	RecommendationsReviewed int             `json:"recommendations_reviewed"`
	OpenedAt                time.Time       `json:"opened_at,omitzero"`
}

type caseHistoryResponse struct {
	Visits       []historyVisitResponse       `json:"visits"`
	Appointments []historyAppointmentResponse `json:"appointments"`
}

type historyVisitResponse struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	VisitedAt   time.Time `json:"visited_at"`
}

type historyAppointmentResponse struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type reviewRecommendationRequest struct {
	Mark          string  `json:"mark"`
	RejectReason  string  `json:"reject_reason"`
	RejectComment string  `json:"reject_comment"`
	PatientText   *string `json:"patient_text"`
}

type addRecommendationRequest struct {
	ServiceCode string `json:"service_code"`
	ServiceName string `json:"service_name"`
	Rationale   string `json:"rationale"`
	PatientText string `json:"patient_text"`
	Mark        string `json:"mark"`
}

type changeUrgencyRequest struct {
	Urgency string `json:"urgency"`
	Reason  string `json:"reason"`
}

func (r reviewRecommendationRequest) toCommand(caseID, recommendationID uuid.UUID, actor string) review.ReviewRecommendation {
	return review.ReviewRecommendation{
		CaseID:           caseID,
		RecommendationID: recommendationID,
		Mark:             domain.Mark(r.Mark),
		RejectReason:     domain.RejectReason(r.RejectReason),
		RejectComment:    r.RejectComment,
		PatientText:      r.PatientText,
		Actor:            actor,
	}
}

func (r addRecommendationRequest) toCommand(caseID uuid.UUID, actor string) review.AddRecommendation {
	return review.AddRecommendation{
		CaseID:      caseID,
		Service:     domain.Service{Code: r.ServiceCode, Name: r.ServiceName},
		Rationale:   r.Rationale,
		PatientText: r.PatientText,
		Mark:        domain.Mark(r.Mark),
		Actor:       actor,
	}
}

func (r changeUrgencyRequest) toCommand(caseID uuid.UUID, actor string) review.ChangeUrgency {
	return review.ChangeUrgency{
		CaseID:  caseID,
		Urgency: domain.Urgency(r.Urgency),
		Reason:  r.Reason,
		Actor:   actor,
	}
}

func newCaseResponse(d cases.Details) caseResponse {
	c := d.Case
	resp := caseResponse{
		ID:                c.ID,
		PatientID:         c.Patient.ID,
		SourceSystem:      c.Patient.SourceSystem,
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
		AssessedAt:        c.AssessedAt,
		UpdatedAt:         c.UpdatedAt,
		Patient:           newPatientResponse(c.Patient),
		Recommendations:   make([]recommendationResponse, 0, len(d.Recommendations)),
	}
	for _, rec := range d.Recommendations {
		resp.Recommendations = append(resp.Recommendations, newRecommendationResponse(rec))
	}
	return resp
}

func newRecommendationResponse(r domain.Recommendation) recommendationResponse {
	resp := recommendationResponse{
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
	}
	if r.Reviewed() {
		resp.Review = &reviewResponse{
			Mark:          string(r.Review.Mark),
			RejectReason:  string(r.Review.RejectReason),
			RejectComment: r.Review.RejectComment,
			ReviewedBy:    r.Review.ReviewedBy,
			ReviewedAt:    r.Review.ReviewedAt,
		}
	}
	return resp
}

func newPatientResponse(p domain.Patient) patientResponse {
	return patientResponse{ID: p.ExternalID, BirthDate: p.BirthDate.Format(time.DateOnly), Sex: string(p.Sex)}
}

func newReviewQueueResponse(items []domain.ReviewQueueItem) reviewQueueResponse {
	resp := reviewQueueResponse{Cases: make([]queueCaseResponse, 0, len(items))}
	for _, item := range items {
		resp.Cases = append(resp.Cases, queueCaseResponse{
			CaseID:                  item.CaseID,
			Urgency:                 string(item.Urgency),
			Modality:                string(item.Modality),
			PerformedAt:             item.PerformedAt,
			ReceivedAt:              item.ReceivedAt,
			Patient:                 newPatientResponse(item.Patient),
			RecommendationsTotal:    item.RecommendationsTotal,
			RecommendationsReviewed: item.RecommendationsReviewed,
			OpenedAt:                item.OpenedAt,
		})
	}
	return resp
}

func newCaseHistoryResponse(h gatewayapi.History) caseHistoryResponse {
	resp := caseHistoryResponse{
		Visits:       make([]historyVisitResponse, 0, len(h.Visits)),
		Appointments: make([]historyAppointmentResponse, 0, len(h.Appointments)),
	}
	for _, v := range h.Visits {
		resp.Visits = append(resp.Visits, historyVisitResponse{ServiceCode: v.ServiceCode, ServiceName: v.ServiceName, VisitedAt: v.VisitedAt})
	}
	for _, a := range h.Appointments {
		resp.Appointments = append(resp.Appointments, historyAppointmentResponse{ServiceCode: a.ServiceCode, ServiceName: a.ServiceName, ScheduledAt: a.ScheduledAt})
	}
	return resp
}
