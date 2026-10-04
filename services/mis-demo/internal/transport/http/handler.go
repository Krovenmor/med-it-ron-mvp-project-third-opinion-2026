package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
)

type ReportSender interface {
	SendReport(ctx context.Context, r demo.Report) (b2b.Delivery, error)
}

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

func (h *Handler) patientHistory(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, demo.HistoryOf(r.PathValue("id")))
}

func (h *Handler) sendReports(w http.ResponseWriter, r *http.Request) {
	reports := demo.Reports()
	deliveries := make([]b2b.Delivery, 0, len(reports))
	for _, report := range reports {
		d, err := h.sender.SendReport(r.Context(), report)
		if err != nil {
			h.log.Error("send report to b2b", zap.Error(err))
			h.writeJSON(w, http.StatusBadGateway, deliveriesResponse{Delivered: deliveries, Error: err.Error()})
			return
		}
		deliveries = append(deliveries, d)
	}
	h.writeJSON(w, http.StatusOK, deliveriesResponse{Delivered: deliveries})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.log.Warn("write response", zap.Error(err))
	}
}

func (h *Handler) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		h.log.Info("request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
