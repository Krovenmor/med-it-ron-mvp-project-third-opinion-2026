//go:build e2e

package e2e

import (
	"time"

	"github.com/google/uuid"
)

type planDTO struct {
	Now     time.Time      `json:"now"`
	Patient planPatientDTO `json:"patient"`
	Cases   []planCaseDTO  `json:"cases"`
}

type planPatientDTO struct {
	FullName string `json:"full_name"`
}

type planCaseDTO struct {
	CaseID          string                  `json:"case_id"`
	Status          string                  `json:"status"`
	UrgentContact   bool                    `json:"urgent_contact"`
	NoFindings      bool                    `json:"no_findings"`
	ConfirmedAt     *time.Time              `json:"confirmed_at"`
	BookBy          *time.Time              `json:"book_by"`
	Study           planStudyDTO            `json:"study"`
	Recommendations []planRecommendationDTO `json:"recommendations"`
}

type planStudyDTO struct {
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
}

type planRecommendationDTO struct {
	ID            string `json:"id"`
	ServiceCode   string `json:"service_code"`
	ServiceName   string `json:"service_name"`
	PatientText   string `json:"patient_text"`
	Mark          string `json:"mark"`
	Declined      bool   `json:"declined"`
	DeclineReason string `json:"decline_reason,omitempty"`
}

type b2bDeclineCall struct {
	RecommendationID string
	Reason           string `json:"reason"`
	Comment          string `json:"comment"`
}

type b2bBookingCall struct {
	RecommendationID string    `json:"recommendation_id"`
	AppointmentID    string    `json:"appointment_id"`
	ScheduledAt      time.Time `json:"scheduled_at"`
	Channel          string    `json:"channel"`
}

type misHistory struct {
	Visits       []misVisit       `json:"visits"`
	Appointments []misAppointment `json:"appointments"`
}

type misVisit struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	VisitedAt   time.Time `json:"visited_at"`
}

type misAppointment struct {
	ID          string    `json:"id"`
	PatientID   string    `json:"patient_id"`
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
	ReferralID  string    `json:"referral_id,omitempty"`
}

type misSlots struct {
	Slots []misSlot `json:"slots"`
}

type misSlot struct {
	ID          string    `json:"id"`
	ServiceCode string    `json:"service_code"`
	StartsAt    time.Time `json:"starts_at"`
}

type misBookingCall struct {
	PatientID   string `json:"patient_id"`
	SlotID      string `json:"slot_id"`
	ServiceName string `json:"service_name"`
	ReferralID  string `json:"referral_id"`
}

type routeView struct {
	Now           time.Time         `json:"now"`
	Patient       routePatientView  `json:"patient"`
	Clinic        clinicView        `json:"clinic"`
	UrgentContact bool              `json:"urgent_contact"`
	PendingReview bool              `json:"pending_review"`
	NoFindings    bool              `json:"no_findings"`
	Progress      *int              `json:"progress"`
	NextStepID    string            `json:"next_step_id"`
	Trunk         []trunkNodeView   `json:"trunk"`
	Steps         []stepView        `json:"steps"`
	Appointments  []appointmentView `json:"appointments"`
}

type routePatientView struct {
	FullName string `json:"full_name"`
}

type clinicView struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type trunkNodeView struct {
	ID       string    `json:"id"`
	Kind     string    `json:"kind"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	AIReport bool      `json:"ai_report"`
	Pending  bool      `json:"pending"`
}

type stepView struct {
	ID            string           `json:"id"`
	Kind          string           `json:"kind"`
	OriginID      string           `json:"origin_id"`
	CaseID        string           `json:"case_id"`
	Title         string           `json:"title"`
	Text          string           `json:"text"`
	Mark          string           `json:"mark"`
	Status        string           `json:"status"`
	Late          bool             `json:"late"`
	AssignedAt    time.Time        `json:"assigned_at"`
	DueAt         time.Time        `json:"due_at"`
	Date          time.Time        `json:"date"`
	InClinic      bool             `json:"in_clinic"`
	Bookable      bool             `json:"bookable"`
	Appointment   *appointmentView `json:"appointment"`
	VisitedAt     time.Time        `json:"visited_at"`
	DeclineReason string           `json:"decline_reason"`
}

type declineRequest struct {
	Reason  string `json:"reason"`
	Comment string `json:"comment,omitempty"`
}

type slotsView struct {
	Slots []slotView `json:"slots"`
}

type slotView struct {
	ID       string    `json:"id"`
	StartsAt time.Time `json:"starts_at"`
}

type bookRequest struct {
	SlotID string `json:"slot_id"`
}

type appointmentView struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

func newRecommendation(serviceCode, serviceName string) planRecommendationDTO {
	return planRecommendationDTO{
		ID:          uuid.NewString(),
		ServiceCode: serviceCode,
		ServiceName: serviceName,
		PatientText: "Текст для пациента: " + serviceName,
		Mark:        "critical",
	}
}

func newConfirmedCase(modality string, performedAt time.Time, recs ...planRecommendationDTO) planCaseDTO {
	confirmedAt := performedAt.Add(6 * time.Hour)
	bookBy := confirmedAt.AddDate(0, 0, 7)
	return planCaseDTO{
		CaseID:          uuid.NewString(),
		Status:          "confirmed",
		ConfirmedAt:     &confirmedAt,
		BookBy:          &bookBy,
		Study:           planStudyDTO{Modality: modality, PerformedAt: performedAt},
		Recommendations: recs,
	}
}

func uniqueID(prefix string) string {
	return prefix + "-" + uuid.NewString()
}
