package http

import (
	"net/http"
	"time"
)

func (h *Handler) advanceClock(w http.ResponseWriter, r *http.Request) {
	var req advanceClockRequest
	if !h.decode(w, r, &req) {
		return
	}
	d, err := time.ParseDuration(req.By)
	if err != nil || d <= 0 {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: `"by" must be a positive duration, e.g. "30m" or "24h"`})
		return
	}
	h.writeJSON(w, http.StatusOK, clockResponse{Now: h.clock.Advance(d)})
}
