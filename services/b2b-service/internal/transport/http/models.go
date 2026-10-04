package http

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/booking"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/cases"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/service/review"
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
	Patient           patientResponse          `json:"patient"`
	Recommendations   []recommendationResponse `json:"recommendations"`
}

type patientResponse struct {
	FullName  string `json:"full_name"`
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
	return patientResponse{
		FullName:  p.FullName,
		BirthDate: p.BirthDate.Format(time.DateOnly),
		Sex:       string(p.Sex),
	}
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
		})
	}
	return resp
}

func newCaseResponse(d cases.Details) caseResponse {
	c := d.Case
	recs := make([]recommendationResponse, 0, len(d.Recommendations))
	for _, r := range d.Recommendations {
		recs = append(recs, newRecommendationResponse(r))
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
		Patient:           newPatientResponse(d.Patient),
		Recommendations:   recs,
	}
}

type bookingRequest struct {
	RecommendationID string `json:"recommendation_id"`
	AppointmentID    string `json:"appointment_id"`
	ScheduledAt      string `json:"scheduled_at"`
	Channel          string `json:"channel"`
}

type bookingResponse struct {
	BookingID        uuid.UUID `json:"booking_id"`
	CaseID           uuid.UUID `json:"case_id"`
	RecommendationID uuid.UUID `json:"recommendation_id"`
	AppointmentID    string    `json:"appointment_id"`
	CaseStatus       string    `json:"case_status"`
}

type patientPlanResponse struct {
	Patient planPatientResponse `json:"patient"`
	Cases   []planCaseResponse  `json:"cases"`
}

type planPatientResponse struct {
	FullName string `json:"full_name"`
}

type planCaseResponse struct {
	CaseID          uuid.UUID                    `json:"case_id"`
	Status          string                       `json:"status"`
	UrgentContact   bool                         `json:"urgent_contact"`
	Study           planStudyResponse            `json:"study"`
	Recommendations []planRecommendationResponse `json:"recommendations"`
}

type planStudyResponse struct {
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
}

type planRecommendationResponse struct {
	ID          uuid.UUID `json:"id"`
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	PatientText string    `json:"patient_text"`
	Mark        string    `json:"mark"`
}

func (r bookingRequest) toCommand(caseID uuid.UUID) (booking.Book, error) {
	recommendationID, err := uuid.Parse(r.RecommendationID)
	if err != nil {
		return booking.Book{}, fmt.Errorf("%w: recommendation_id must be a UUID", domain.ErrInvalidInput)
	}
	scheduledAt, err := parseOptionalTime("scheduled_at", r.ScheduledAt, time.RFC3339, "an RFC 3339 date-time")
	if err != nil {
		return booking.Book{}, err
	}
	return booking.Book{
		CaseID:           caseID,
		RecommendationID: recommendationID,
		AppointmentID:    r.AppointmentID,
		ScheduledAt:      scheduledAt,
		Channel:          domain.BookingChannel(r.Channel),
	}, nil
}

func newBookingResponse(r booking.Result) bookingResponse {
	return bookingResponse{
		BookingID:        r.Booking.ID,
		CaseID:           r.Booking.CaseID,
		RecommendationID: r.Booking.RecommendationID,
		AppointmentID:    r.Booking.AppointmentID,
		CaseStatus:       string(r.CaseStatus),
	}
}

func newPatientPlanResponse(p domain.PatientPlan) patientPlanResponse {
	resp := patientPlanResponse{
		Patient: planPatientResponse{FullName: p.Patient.FullName},
		Cases:   make([]planCaseResponse, 0, len(p.Cases)),
	}
	for _, c := range p.Cases {
		resp.Cases = append(resp.Cases, newPlanCaseResponse(c))
	}
	return resp
}

func newPlanCaseResponse(c domain.PlanCase) planCaseResponse {
	recs := make([]planRecommendationResponse, 0, len(c.Recommendations))
	for _, r := range c.Recommendations {
		recs = append(recs, planRecommendationResponse{
			ID:          r.ID,
			ServiceCode: r.ServiceCode,
			ServiceName: r.ServiceName,
			PatientText: r.PatientText,
			Mark:        string(r.Review.Mark),
		})
	}
	return planCaseResponse{
		CaseID:        c.Case.ID,
		Status:        string(c.Case.Status),
		UrgentContact: c.UrgentContact,
		Study: planStudyResponse{
			Modality:    string(c.Case.Study.Modality),
			PerformedAt: c.Case.Study.PerformedAt,
		},
		Recommendations: recs,
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
