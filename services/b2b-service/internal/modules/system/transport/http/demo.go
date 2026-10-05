package http

import (
	"net/http"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/httpx"
)

func (h *Handler) advanceClock(w http.ResponseWriter, r *http.Request) {
	var req advanceClockRequest
	if !h.Decode(w, r, &req) {
		return
	}
	d, err := time.ParseDuration(req.By)
	if err != nil || d <= 0 {
		h.WriteJSON(w, http.StatusBadRequest, httpx.ErrorResponse{Error: `"by" must be a positive duration, e.g. "30m" or "24h"`})
		return
	}
	h.WriteJSON(w, http.StatusOK, clockResponse{Now: h.clock.Advance(d)})
}

func (h *Handler) resetDemo(w http.ResponseWriter, r *http.Request) {
	if err := h.demo.Reset(r.Context()); err != nil {
		h.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
