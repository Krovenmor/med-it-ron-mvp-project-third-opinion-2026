//go:build e2e

package e2e

import (
	"time"

	"github.com/google/uuid"
)

type planDTO struct {
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
	Study           planStudyDTO            `json:"study"`
	Recommendations []planRecommendationDTO `json:"recommendations"`
}

type planStudyDTO struct {
	Modality    string    `json:"modality"`
	PerformedAt time.Time `json:"performed_at"`
}

type planRecommendationDTO struct {
	ID          string `json:"id"`
	ServiceCode string `json:"service_code"`
	ServiceName string `json:"service_name"`
	PatientText string `json:"patient_text"`
	Mark        string `json:"mark"`
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

type graphView struct {
	Patient     graphPatientView `json:"patient"`
	ClinicPhone string           `json:"clinic_phone"`
	Nodes       []nodeView       `json:"nodes"`
	Edges       []edgeView       `json:"edges"`
}

type graphPatientView struct {
	FullName string `json:"full_name"`
}

type nodeView struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	Date     time.Time `json:"date"`
	Status   string    `json:"status"`
	Mark     string    `json:"mark"`
	InClinic *bool     `json:"in_clinic"`
	Bookable *bool     `json:"bookable"`
}

type edgeView struct {
	From string `json:"from"`
	To   string `json:"to"`
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
	return planCaseDTO{
		CaseID:          uuid.NewString(),
		Status:          "confirmed",
		Study:           planStudyDTO{Modality: modality, PerformedAt: performedAt},
		Recommendations: recs,
	}
}

func uniqueID(prefix string) string {
	return prefix + "-" + uuid.NewString()
}
