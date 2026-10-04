package http

import "net/http"

func (h *Handler) listSlots(w http.ResponseWriter, r *http.Request) {
	slots, err := h.bookings.Slots(r.Context(), r.PathValue("patient_id"), r.PathValue("rec_id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusOK, newSlotsResponse(slots))
}

func (h *Handler) bookAppointment(w http.ResponseWriter, r *http.Request) {
	var req bookAppointmentRequest
	if !h.decode(w, r, &req) {
		return
	}
	appointment, err := h.bookings.Book(r.Context(), r.PathValue("patient_id"), r.PathValue("rec_id"), req.SlotID)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, newAppointmentResponse(appointment))
}
