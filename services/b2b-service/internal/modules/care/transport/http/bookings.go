package http

import "net/http"

func (h *Handler) registerBooking(w http.ResponseWriter, r *http.Request) {
	caseID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	var req bookingRequest
	if !h.Decode(w, r, &req) {
		return
	}
	cmd, err := req.toCommand(caseID)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}

	result, err := h.bookings.Book(r.Context(), cmd)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}

	status := http.StatusOK
	if result.Created {
		status = http.StatusCreated
	}
	h.WriteJSON(w, status, newBookingResponse(result))
}
