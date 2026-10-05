package b2b

import (
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type planResponse struct {
	Now     time.Time           `json:"now"`
	Patient planPatientResponse `json:"patient"`
	Cases   []planCaseResponse  `json:"cases"`
}

type planPatientResponse struct {
	FullName string `json:"full_name"`
}

type planCaseResponse struct {
	CaseID          string                       `json:"case_id"`
	Status          string                       `json:"status"`
	UrgentContact   bool                         `json:"urgent_contact"`
	NoFindings      bool                         `json:"no_findings"`
	ConfirmedAt     *time.Time                   `json:"confirmed_at"`
	BookBy          *time.Time                   `json:"book_by"`
	Study           planStudyResponse            `json:"study"`
	Recommendations []planRecommendationResponse `json:"recommendations"`
}

type planStudyResponse struct {
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
}

type planRecommendationResponse struct {
	ID            string `json:"id"`
	ServiceCode   string `json:"service_code"`
	ServiceName   string `json:"service_name"`
	PatientText   string `json:"patient_text"`
	Mark          string `json:"mark"`
	Declined      bool   `json:"declined"`
	DeclineReason string `json:"decline_reason"`
}

type declineRequest struct {
	Reason  string `json:"reason"`
	Comment string `json:"comment"`
}

type bookingRequest struct {
	RecommendationID string    `json:"recommendation_id"`
	AppointmentID    string    `json:"appointment_id"`
	ScheduledAt      time.Time `json:"scheduled_at"`
	Channel          string    `json:"channel"`
}

func (r planResponse) toDomain() domain.Plan {
	plan := domain.Plan{Now: r.Now, PatientName: r.Patient.FullName, Cases: make([]domain.PlanCase, 0, len(r.Cases))}
	for _, c := range r.Cases {
		plan.Cases = append(plan.Cases, c.toDomain())
	}
	return plan
}

func (c planCaseResponse) toDomain() domain.PlanCase {
	recs := make([]domain.Recommendation, 0, len(c.Recommendations))
	for _, r := range c.Recommendations {
		recs = append(recs, r.toDomain(c.CaseID))
	}
	return domain.PlanCase{
		ID:              c.CaseID,
		Status:          c.Status,
		UrgentContact:   c.UrgentContact,
		NoFindings:      c.NoFindings,
		Modality:        c.Study.Modality,
		PerformedAt:     c.Study.PerformedAt,
		ConfirmedAt:     valueOf(c.ConfirmedAt),
		BookBy:          valueOf(c.BookBy),
		Recommendations: recs,
	}
}

func (r planRecommendationResponse) toDomain(caseID string) domain.Recommendation {
	return domain.Recommendation{
		ID:            r.ID,
		CaseID:        caseID,
		ServiceCode:   r.ServiceCode,
		ServiceName:   r.ServiceName,
		PatientText:   r.PatientText,
		Mark:          r.Mark,
		Declined:      r.Declined,
		DeclineReason: r.DeclineReason,
	}
}

func valueOf[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}

func newBookingRequest(b domain.Booking) bookingRequest {
	return bookingRequest{
		RecommendationID: b.RecommendationID,
		AppointmentID:    b.AppointmentID,
		ScheduledAt:      b.ScheduledAt,
		Channel:          "self",
	}
}
