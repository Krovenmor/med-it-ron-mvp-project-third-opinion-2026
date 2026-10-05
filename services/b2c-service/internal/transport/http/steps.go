package http

import "net/http"

func (h *Handler) declineStep(w http.ResponseWriter, r *http.Request) {
	var req declineRequest
	if !h.decode(w, r, &req) {
		return
	}
	err := h.steps.Decline(r.Context(), r.PathValue("patient_id"), req.toDomain(r.PathValue("rec_id")))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) requestHelp(w http.ResponseWriter, r *http.Request) {
	if err := h.steps.RequestHelp(r.Context(), r.PathValue("patient_id"), r.PathValue("rec_id")); err != nil {
		h.writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
