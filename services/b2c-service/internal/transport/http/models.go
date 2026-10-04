package http

import (
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type graphResponse struct {
	Patient     graphPatientResponse `json:"patient"`
	ClinicPhone string               `json:"clinic_phone"`
	Nodes       []nodeResponse       `json:"nodes"`
	Edges       []edgeResponse       `json:"edges"`
}

type graphPatientResponse struct {
	FullName string `json:"full_name"`
}

type nodeResponse struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Text     string    `json:"text,omitempty"`
	Date     time.Time `json:"date,omitzero"`
	Status   string    `json:"status,omitempty"`
	Mark     string    `json:"mark,omitempty"`
	InClinic *bool     `json:"in_clinic,omitempty"`
	Bookable *bool     `json:"bookable,omitempty"`
}

type edgeResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
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

func newGraphResponse(g domain.Graph) graphResponse {
	resp := graphResponse{
		Patient:     graphPatientResponse{FullName: g.PatientName},
		ClinicPhone: g.ClinicPhone,
		Nodes:       make([]nodeResponse, 0, len(g.Nodes)),
		Edges:       make([]edgeResponse, 0, len(g.Edges)),
	}
	for _, n := range g.Nodes {
		resp.Nodes = append(resp.Nodes, newNodeResponse(n))
	}
	for _, e := range g.Edges {
		resp.Edges = append(resp.Edges, edgeResponse{From: e.From, To: e.To})
	}
	return resp
}

func newNodeResponse(n domain.Node) nodeResponse {
	resp := nodeResponse{
		ID:     n.ID,
		Type:   string(n.Type),
		Title:  n.Title,
		Text:   n.Text,
		Date:   n.Date,
		Status: string(n.Status),
		Mark:   n.Mark,
	}
	if n.Type == domain.NodeTypeRecommendation {
		resp.InClinic = &n.InClinic
		resp.Bookable = &n.Bookable
	}
	return resp
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
