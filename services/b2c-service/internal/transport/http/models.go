package http

import (
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type routeResponse struct {
	Now           time.Time             `json:"now"`
	Patient       routePatientResponse  `json:"patient"`
	Clinic        clinicResponse        `json:"clinic"`
	UrgentContact bool                  `json:"urgent_contact"`
	PendingReview bool                  `json:"pending_review"`
	NoFindings    bool                  `json:"no_findings"`
	Progress      *int                  `json:"progress"`
	NextStepID    string                `json:"next_step_id"`
	Trunk         []trunkNodeResponse   `json:"trunk"`
	Steps         []stepResponse        `json:"steps"`
	Appointments  []appointmentResponse `json:"appointments"`
}

type routePatientResponse struct {
	FullName string `json:"full_name"`
}

type clinicResponse struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type trunkNodeResponse struct {
	ID       string    `json:"id"`
	Kind     string    `json:"kind"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	AIReport bool      `json:"ai_report"`
	Pending  bool      `json:"pending"`
}

type stepResponse struct {
	ID            string               `json:"id"`
	Kind          string               `json:"kind"`
	OriginID      string               `json:"origin_id"`
	CaseID        string               `json:"case_id,omitempty"`
	Title         string               `json:"title"`
	Text          string               `json:"text,omitempty"`
	Mark          string               `json:"mark,omitempty"`
	Status        string               `json:"status"`
	Late          bool                 `json:"late"`
	AssignedAt    time.Time            `json:"assigned_at,omitzero"`
	DueAt         time.Time            `json:"due_at,omitzero"`
	Date          time.Time            `json:"date"`
	InClinic      bool                 `json:"in_clinic"`
	Bookable      bool                 `json:"bookable"`
	Appointment   *appointmentResponse `json:"appointment"`
	VisitedAt     time.Time            `json:"visited_at,omitzero"`
	DeclineReason string               `json:"decline_reason,omitempty"`
}

type declineRequest struct {
	Reason  string `json:"reason"`
	Comment string `json:"comment"`
}

type slotsResponse struct {
	Slots []slotResponse `json:"slots"`
}

type slotResponse struct {
	ID       string    `json:"id"`
	StartsAt time.Time `json:"starts_at"`
}

type bookAppointmentRequest struct {
	SlotID string `json:"slot_id"`
}

type appointmentResponse struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func newRouteResponse(r domain.Route) routeResponse {
	resp := routeResponse{
		Now:           r.Now,
		Patient:       routePatientResponse{FullName: r.PatientName},
		Clinic:        clinicResponse{Name: r.Clinic.Name, Phone: r.Clinic.Phone, Address: r.Clinic.Address},
		UrgentContact: r.UrgentContact,
		PendingReview: r.PendingReview,
		NoFindings:    r.NoFindings,
		Progress:      r.Progress,
		NextStepID:    r.NextStepID,
		Trunk:         make([]trunkNodeResponse, 0, len(r.Trunk)),
		Steps:         make([]stepResponse, 0, len(r.Steps)),
		Appointments:  make([]appointmentResponse, 0, len(r.Appointments)),
	}
	for _, n := range r.Trunk {
		resp.Trunk = append(resp.Trunk, trunkNodeResponse{
			ID:       n.ID,
			Kind:     string(n.Kind),
			Title:    n.Title,
			Date:     n.Date,
			AIReport: n.AIReport,
			Pending:  n.Pending,
		})
	}
	for _, s := range r.Steps {
		resp.Steps = append(resp.Steps, newStepResponse(s))
	}
	for _, a := range r.Appointments {
		resp.Appointments = append(resp.Appointments, newAppointmentResponse(a))
	}
	return resp
}

func newStepResponse(s domain.Step) stepResponse {
	resp := stepResponse{
		ID:            s.ID,
		Kind:          string(s.Kind),
		OriginID:      s.OriginID,
		CaseID:        s.CaseID,
		Title:         s.Title,
		Text:          s.Text,
		Mark:          s.Mark,
		Status:        string(s.Status),
		Late:          s.Late,
		AssignedAt:    s.AssignedAt,
		DueAt:         s.DueAt,
		Date:          s.Date,
		InClinic:      s.InClinic,
		Bookable:      s.Bookable,
		VisitedAt:     s.VisitedAt,
		DeclineReason: s.DeclineReason,
	}
	if s.Appointment != nil {
		a := newAppointmentResponse(*s.Appointment)
		resp.Appointment = &a
	}
	return resp
}

func (r declineRequest) toDomain(recommendationID string) domain.Decline {
	return domain.Decline{RecommendationID: recommendationID, Reason: r.Reason, Comment: r.Comment}
}

func newSlotsResponse(slots []domain.Slot) slotsResponse {
	resp := slotsResponse{Slots: make([]slotResponse, 0, len(slots))}
	for _, s := range slots {
		resp.Slots = append(resp.Slots, slotResponse{ID: s.ID, StartsAt: s.StartsAt})
	}
	return resp
}

func newAppointmentResponse(a domain.Appointment) appointmentResponse {
	return appointmentResponse{ID: a.ID, ServiceName: a.ServiceName, ScheduledAt: a.ScheduledAt}
}
