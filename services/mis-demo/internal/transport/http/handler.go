package http

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	sender    ReportSender
	scheduler Scheduler
	log       *zap.Logger
}

func NewHandler(sender ReportSender, scheduler Scheduler, log *zap.Logger) *Handler {
	return &Handler{sender: sender, scheduler: scheduler, log: log.Named("http")}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/patients/{id}/history", h.patientHistory)
	mux.HandleFunc("GET /api/v1/slots", h.listSlots)
	mux.HandleFunc("POST /api/v1/appointments", h.bookAppointment)
	mux.HandleFunc("GET /api/v1/services", h.searchServices)
	mux.HandleFunc("POST /demo/reports", h.sendReports)
	mux.HandleFunc("POST /demo/reset", h.resetDemo)
	return h.logRequests(mux)
}
