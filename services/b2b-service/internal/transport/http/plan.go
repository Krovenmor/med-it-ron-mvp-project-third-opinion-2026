package http

import "net/http"

func (h *Handler) patientPlan(w http.ResponseWriter, r *http.Request) {
	sourceSystem := r.URL.Query().Get("source_system")
	if sourceSystem == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "source_system query parameter is required"})
		return
	}
	plan, err := h.plan.Get(r.Context(), sourceSystem, r.PathValue("external_id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusOK, newPatientPlanResponse(plan))
}
