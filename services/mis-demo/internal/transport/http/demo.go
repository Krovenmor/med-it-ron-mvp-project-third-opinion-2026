package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/b2b"
)

func (h *Handler) sendReports(w http.ResponseWriter, r *http.Request) {
	reports := h.scenario.Reports()
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

func (h *Handler) resetDemo(w http.ResponseWriter, _ *http.Request) {
	h.scheduler.Reset()
	h.clock.Reset()
	h.scenario.Reset(h.clock.Now())
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) advanceClock(w http.ResponseWriter, r *http.Request) {
	var req advanceClockRequest
	if err := decodeJSON(r, &req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "malformed JSON body"})
		return
	}
	d, err := time.ParseDuration(req.By)
	if err != nil || d <= 0 {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: `"by" must be a positive duration, e.g. "30m" or "24h"`})
		return
	}
	h.writeJSON(w, http.StatusOK, clockResponse{Now: h.clock.Advance(d)})
}
