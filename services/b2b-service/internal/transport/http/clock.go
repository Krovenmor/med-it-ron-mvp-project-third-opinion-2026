package http

import "net/http"

func (h *Handler) currentTime(w http.ResponseWriter, _ *http.Request) {
	h.writeJSON(w, http.StatusOK, clockResponse{Now: h.clock.Now()})
}
