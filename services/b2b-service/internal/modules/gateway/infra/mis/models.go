package mis

import (
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
)

type historyResponse struct {
	Visits       []visitDTO       `json:"visits"`
	Appointments []appointmentDTO `json:"appointments"`
}

type visitDTO struct {
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	VisitedAt   time.Time `json:"visited_at"`
}

type appointmentDTO struct {
	ID          string    `json:"id"`
	ServiceCode string    `json:"service_code"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
	ReferralID  string    `json:"referral_id"`
}

type slotsResponse struct {
	Slots []slotDTO `json:"slots"`
}

type slotDTO struct {
	ID          string    `json:"id"`
	ServiceCode string    `json:"service_code"`
	StartsAt    time.Time `json:"starts_at"`
}

type appointmentRequest struct {
	PatientID   string `json:"patient_id"`
	SlotID      string `json:"slot_id"`
	ServiceName string `json:"service_name"`
	ReferralID  string `json:"referral_id"`
}

func (r historyResponse) toDomain() domain.PatientHistory {
	h := domain.PatientHistory{
		Visits:       make([]domain.Visit, 0, len(r.Visits)),
		Appointments: make([]domain.Appointment, 0, len(r.Appointments)),
	}
	for _, v := range r.Visits {
		h.Visits = append(h.Visits, v.toDomain())
	}
	for _, a := range r.Appointments {
		h.Appointments = append(h.Appointments, a.toDomain())
	}
	return h
}

func (v visitDTO) toDomain() domain.Visit {
	return domain.Visit{ServiceCode: v.ServiceCode, ServiceName: v.ServiceName, VisitedAt: v.VisitedAt}
}

func (a appointmentDTO) toDomain() domain.Appointment {
	return domain.Appointment{
		ID:          a.ID,
		ServiceCode: a.ServiceCode,
		ServiceName: a.ServiceName,
		ScheduledAt: a.ScheduledAt,
		ReferralID:  a.ReferralID,
	}
}

func (r slotsResponse) toDomain() []domain.Slot {
	slots := make([]domain.Slot, 0, len(r.Slots))
	for _, s := range r.Slots {
		slots = append(slots, s.toDomain())
	}
	return slots
}

func (s slotDTO) toDomain() domain.Slot {
	return domain.Slot{ID: s.ID, ServiceCode: s.ServiceCode, StartsAt: s.StartsAt}
}

func newAppointmentRequest(r domain.AppointmentRequest) appointmentRequest {
	return appointmentRequest{
		PatientID:   r.PatientID,
		SlotID:      r.SlotID,
		ServiceName: r.ServiceName,
		ReferralID:  r.ReferralID,
	}
}
