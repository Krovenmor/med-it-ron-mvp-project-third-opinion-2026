package http

import "net/http"

func (h *Handler) registerBooking(w http.ResponseWriter, r *http.Request) {
	caseID, ok := h.pathID(w, r, "id")
	if !ok {
		return
	}
	var req bookingRequest
	if !h.decode(w, r, &req) {
		return
	}
	cmd, err := req.toCommand(caseID)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	result, err := h.bookings.Book(r.Context(), cmd)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	status := http.StatusOK
	if result.Created {
		status = http.StatusCreated
	}
	h.writeJSON(w, status, newBookingResponse(result))
}
