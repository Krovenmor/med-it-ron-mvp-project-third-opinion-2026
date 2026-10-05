package http

import "net/http"

func (h *Handler) declineRecommendation(w http.ResponseWriter, r *http.Request) {
	caseID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	recID, ok := h.PathID(w, r, "rec_id")
	if !ok {
		return
	}
	var req patientDeclineRequest
	if !h.Decode(w, r, &req) {
		return
	}
	result, err := h.patient.Decline(r.Context(), req.toCommand(caseID, recID))
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, newPatientDeclineResponse(result))
}

func (h *Handler) requestHelp(w http.ResponseWriter, r *http.Request) {
	caseID, ok := h.PathID(w, r, "id")
	if !ok {
		return
	}
	result, err := h.patient.RequestHelp(r.Context(), caseID)
	if err != nil {
		h.WriteError(w, r, err)
		return
	}
	status := http.StatusOK
	if result.Created {
		status = http.StatusCreated
	}
	h.WriteJSON(w, status, helpRequestResponse{TaskID: result.Task.ID, Created: result.Created})
}
