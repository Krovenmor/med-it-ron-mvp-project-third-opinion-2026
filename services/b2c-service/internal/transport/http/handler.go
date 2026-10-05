package http

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	routes   Routes
	bookings Bookings
	steps    Steps
	log      *zap.Logger
}

func NewHandler(routes Routes, bookings Bookings, steps Steps, log *zap.Logger) *Handler {
	return &Handler{routes: routes, bookings: bookings, steps: steps, log: log.Named("http")}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/patients/{patient_id}/route", h.patientRoute)
	mux.HandleFunc("GET /api/v1/patients/{patient_id}/recommendations/{rec_id}/slots", h.listSlots)
	mux.HandleFunc("POST /api/v1/patients/{patient_id}/recommendations/{rec_id}/appointments", h.bookAppointment)
	mux.HandleFunc("POST /api/v1/patients/{patient_id}/recommendations/{rec_id}/decline", h.declineStep)
	mux.HandleFunc("POST /api/v1/patients/{patient_id}/recommendations/{rec_id}/help-requests", h.requestHelp)
	return h.recoverPanics(h.logRequests(mux))
}
