package http

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	sender ReportSender
	log    *zap.Logger
}

func NewHandler(sender ReportSender, log *zap.Logger) *Handler {
	return &Handler{sender: sender, log: log.Named("http")}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/patients/{id}/history", h.patientHistory)
	mux.HandleFunc("POST /demo/reports", h.sendReports)
	return h.logRequests(mux)
}
