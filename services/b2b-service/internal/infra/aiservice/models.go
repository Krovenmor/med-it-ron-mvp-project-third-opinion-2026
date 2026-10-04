package aiservice

import (
	"time"

	"github.com/google/uuid"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type assessmentRequest struct {
	CaseID                 uuid.UUID        `json:"case_id"`
	Study                  studyDTO         `json:"study"`
	Conclusion             string           `json:"conclusion"`
	Patient                patientDTO       `json:"patient"`
	CurrentRecommendations []serviceDTO     `json:"current_recommendations"`
	Visits                 []visitDTO       `json:"visits"`
	Appointments           []appointmentDTO `json:"appointments"`
}

type studyDTO struct {
	Modality    string    `json:"modality"`
	BodySite    string    `json:"body_site"`
	PerformedAt time.Time `json:"performed_at"`
}

type patientDTO struct {
	Age int    `json:"age"`
	Sex string `json:"sex"`
}

type serviceDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type visitDTO struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	VisitedAt   time.Time `json:"visited_at"`
}

type appointmentDTO struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

func newAssessmentRequest(req domain.AssessmentRequest) assessmentRequest {
	out := assessmentRequest{
		CaseID: req.CaseID,
		Study: studyDTO{
			Modality:    string(req.Modality),
			BodySite:    req.BodySite,
			PerformedAt: req.PerformedAt,
		},
		Conclusion:             req.Conclusion,
		Patient:                patientDTO{Age: req.PatientAge, Sex: string(req.PatientSex)},
		CurrentRecommendations: make([]serviceDTO, 0, len(req.CurrentRecommendations)),
		Visits:                 make([]visitDTO, 0, len(req.History.Visits)),
		Appointments:           make([]appointmentDTO, 0, len(req.History.Appointments)),
	}
	for _, s := range req.CurrentRecommendations {
		out.CurrentRecommendations = append(out.CurrentRecommendations, serviceDTO{Code: s.Code, Name: s.Name})
	}
	for _, v := range req.History.Visits {
		out.Visits = append(out.Visits, visitDTO{ServiceCode: v.ServiceCode, ServiceName: v.ServiceName, VisitedAt: v.VisitedAt})
	}
	for _, a := range req.History.Appointments {
		out.Appointments = append(out.Appointments, appointmentDTO{ServiceCode: a.ServiceCode, ServiceName: a.ServiceName, ScheduledAt: a.ScheduledAt})
	}
	return out
}

type assessmentResponse struct {
	Urgency           string              `json:"urgency"`
	CatalogVersion    string              `json:"catalog_version"`
	GuidelinesVersion string              `json:"guidelines_version"`
	Recommendations   []recommendationDTO `json:"recommendations"`
}

type recommendationDTO struct {
	ServiceCode   string `json:"service_code"`
	ServiceName   string `json:"service_name"`
	Importance    string `json:"importance"`
	Rationale     string `json:"rationale"`
	GuidelineRef  string `json:"guideline_ref"`
	PatientText   string `json:"patient_text"`
	AlreadyBooked bool   `json:"already_booked"`
}

func (r assessmentResponse) toDomain() domain.Assessment {
	recs := make([]domain.Recommendation, 0, len(r.Recommendations))
	for _, rec := range r.Recommendations {
		recs = append(recs, domain.Recommendation{
			ServiceCode:   rec.ServiceCode,
			ServiceName:   rec.ServiceName,
			Importance:    domain.Importance(rec.Importance),
			Rationale:     rec.Rationale,
			GuidelineRef:  rec.GuidelineRef,
			PatientText:   rec.PatientText,
			AlreadyBooked: rec.AlreadyBooked,
		})
	}
	return domain.Assessment{
		Urgency:           domain.Urgency(r.Urgency),
		CatalogVersion:    r.CatalogVersion,
		GuidelinesVersion: r.GuidelinesVersion,
		Recommendations:   recs,
	}
}
