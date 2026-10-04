package http

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	graphs   Graphs
	bookings Bookings
	log      *zap.Logger
}

func NewHandler(graphs Graphs, bookings Bookings, log *zap.Logger) *Handler {
	return &Handler{graphs: graphs, bookings: bookings, log: log.Named("http")}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/patients/{patient_id}/graph", h.patientGraph)
	mux.HandleFunc("GET /api/v1/patients/{patient_id}/recommendations/{rec_id}/slots", h.listSlots)
	mux.HandleFunc("POST /api/v1/patients/{patient_id}/recommendations/{rec_id}/appointments", h.bookAppointment)
	return h.recoverPanics(h.logRequests(mux))
}
