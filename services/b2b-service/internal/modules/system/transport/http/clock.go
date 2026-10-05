package http

import "net/http"

func (h *Handler) currentTime(w http.ResponseWriter, _ *http.Request) {
	h.WriteJSON(w, http.StatusOK, clockResponse{Now: h.clock.Now()})
}
