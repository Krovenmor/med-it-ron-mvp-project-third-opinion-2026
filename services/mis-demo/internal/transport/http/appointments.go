package http

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) listSlots(w http.ResponseWriter, r *http.Request) {
	serviceCode := r.URL.Query().Get("service_code")
	if serviceCode == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "service_code query parameter is required"})
		return
	}
	from, ok := slotsFrom(r.URL.Query().Get("from"), h.clock.Now())
	if !ok {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "from must be an RFC 3339 date-time"})
		return
	}
	h.writeJSON(w, http.StatusOK, slotsResponse{Slots: h.scheduler.Slots(serviceCode, from)})
}

func (h *Handler) bookAppointment(w http.ResponseWriter, r *http.Request) {
	var req bookAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "malformed JSON body"})
		return
	}
	appointment, created, err := h.scheduler.Book(req.toRequest(), h.clock.Now())
	if err != nil {
		h.writeError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	h.writeJSON(w, status, appointment)
}
